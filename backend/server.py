from fastapi import FastAPI, APIRouter, HTTPException, BackgroundTasks
from dotenv import load_dotenv
from starlette.middleware.cors import CORSMiddleware
from motor.motor_asyncio import AsyncIOMotorClient
import os
import logging
from pathlib import Path
from pydantic import BaseModel, Field, ConfigDict
from typing import List, Optional, Dict, Any
import uuid
from datetime import datetime, timezone
from enum import Enum
import subprocess
import json
import asyncio

ROOT_DIR = Path(__file__).parent
load_dotenv(ROOT_DIR / '.env')

# MongoDB connection
mongo_url = os.environ['MONGO_URL']
client = AsyncIOMotorClient(mongo_url)
db = client[os.environ.get('DB_NAME', 'linkedin_automation')]

# Create the main app without a prefix
app = FastAPI(title="LinkedIn Automation Dashboard API")

# Create a router with the /api prefix
api_router = APIRouter(prefix="/api")

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# ============== ENUMS ==============
class ConnectionStatus(str, Enum):
    PENDING = "pending"
    ACCEPTED = "accepted"
    DECLINED = "declined"
    FAILED = "failed"

class MessageStatus(str, Enum):
    SENT = "sent"
    DELIVERED = "delivered"
    READ = "read"
    FAILED = "failed"

class AutomationStatus(str, Enum):
    IDLE = "idle"
    RUNNING = "running"
    PAUSED = "paused"
    ERROR = "error"

# ============== MODELS ==============

# Configuration Models
class LinkedInCredentials(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    email: str
    password: str
    encrypted: bool = False
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

class CredentialsCreate(BaseModel):
    email: str
    password: str

class SearchCriteria(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    name: str
    job_titles: List[str] = []
    companies: List[str] = []
    locations: List[str] = []
    keywords: List[str] = []
    is_active: bool = True
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

class SearchCriteriaCreate(BaseModel):
    name: str
    job_titles: List[str] = []
    companies: List[str] = []
    locations: List[str] = []
    keywords: List[str] = []

class MessageTemplate(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    name: str
    template_type: str  # "connection" or "follow_up"
    content: str
    variables: List[str] = []  # e.g., ["firstName", "company", "jobTitle"]
    is_active: bool = True
    usage_count: int = 0
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

class MessageTemplateCreate(BaseModel):
    name: str
    template_type: str
    content: str
    variables: List[str] = []

class RateLimitConfig(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    daily_connection_limit: int = 50
    daily_message_limit: int = 100
    min_action_delay_ms: int = 5000  # 5 seconds
    max_action_delay_ms: int = 15000  # 15 seconds
    business_hours_start: int = 9  # 9 AM
    business_hours_end: int = 18  # 6 PM
    skip_weekends: bool = True
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

class RateLimitConfigUpdate(BaseModel):
    daily_connection_limit: Optional[int] = None
    daily_message_limit: Optional[int] = None
    min_action_delay_ms: Optional[int] = None
    max_action_delay_ms: Optional[int] = None
    business_hours_start: Optional[int] = None
    business_hours_end: Optional[int] = None
    skip_weekends: Optional[bool] = None

# Connection Models
class Connection(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    profile_url: str
    first_name: Optional[str] = None
    last_name: Optional[str] = None
    job_title: Optional[str] = None
    company: Optional[str] = None
    location: Optional[str] = None
    note_sent: Optional[str] = None
    status: ConnectionStatus = ConnectionStatus.PENDING
    search_criteria_id: Optional[str] = None
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    accepted_at: Optional[datetime] = None

# Message Models
class Message(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    connection_id: str
    content: str
    template_id: Optional[str] = None
    status: MessageStatus = MessageStatus.SENT
    sent_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

# Activity Log Models
class ActivityLog(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    action_type: str  # "login", "search", "connect", "message", "error"
    description: str
    details: Optional[Dict[str, Any]] = None
    status: str = "success"  # "success", "failure", "warning"
    timestamp: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

# Automation State
class AutomationState(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = "automation_state"
    status: AutomationStatus = AutomationStatus.IDLE
    current_task: Optional[str] = None
    connections_today: int = 0
    messages_today: int = 0
    last_action_at: Optional[datetime] = None
    started_at: Optional[datetime] = None
    error_message: Optional[str] = None

# Dashboard Stats
class DashboardStats(BaseModel):
    total_connections: int = 0
    pending_connections: int = 0
    accepted_connections: int = 0
    total_messages: int = 0
    messages_today: int = 0
    connections_today: int = 0
    automation_status: AutomationStatus = AutomationStatus.IDLE
    success_rate: float = 0.0

# Stealth Config
class StealthConfig(BaseModel):
    model_config = ConfigDict(extra="ignore")
    id: str = "stealth_config"
    # Bézier Curve Settings
    bezier_enabled: bool = True
    bezier_overshoot_probability: float = 0.15
    
    # Timing Settings
    typing_min_delay_ms: int = 50
    typing_max_delay_ms: int = 150
    typo_probability: float = 0.05
    
    # Browser Fingerprint
    rotate_user_agent: bool = True
    randomize_viewport: bool = True
    disable_webdriver_flag: bool = True
    
    # Scrolling
    scroll_min_speed: int = 50
    scroll_max_speed: int = 300
    scroll_back_probability: float = 0.1
    
    # Mouse Movement
    hover_before_click: bool = True
    random_cursor_movement: bool = True
    
    # Activity Scheduling
    respect_business_hours: bool = True
    include_break_patterns: bool = True
    
    # Rate Limiting
    enable_token_bucket: bool = True
    cooldown_after_bulk: bool = True
    
    # Headers
    randomize_headers: bool = True
    
    # Network
    simulate_network_latency: bool = True
    
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

class StealthConfigUpdate(BaseModel):
    bezier_enabled: Optional[bool] = None
    bezier_overshoot_probability: Optional[float] = None
    typing_min_delay_ms: Optional[int] = None
    typing_max_delay_ms: Optional[int] = None
    typo_probability: Optional[float] = None
    rotate_user_agent: Optional[bool] = None
    randomize_viewport: Optional[bool] = None
    disable_webdriver_flag: Optional[bool] = None
    scroll_min_speed: Optional[int] = None
    scroll_max_speed: Optional[int] = None
    scroll_back_probability: Optional[float] = None
    hover_before_click: Optional[bool] = None
    random_cursor_movement: Optional[bool] = None
    respect_business_hours: Optional[bool] = None
    include_break_patterns: Optional[bool] = None
    enable_token_bucket: Optional[bool] = None
    cooldown_after_bulk: Optional[bool] = None
    randomize_headers: Optional[bool] = None
    simulate_network_latency: Optional[bool] = None

# Helper function to serialize datetime
def serialize_doc(doc: dict) -> dict:
    for key, value in doc.items():
        if isinstance(value, datetime):
            doc[key] = value.isoformat()
    return doc

def deserialize_doc(doc: dict) -> dict:
    datetime_fields = ['created_at', 'updated_at', 'timestamp', 'sent_at', 'accepted_at', 'started_at', 'last_action_at']
    for field in datetime_fields:
        if field in doc and isinstance(doc[field], str):
            try:
                doc[field] = datetime.fromisoformat(doc[field])
            except:
                pass
    return doc

# ============== CREDENTIALS ENDPOINTS ==============
@api_router.post("/credentials", response_model=dict)
async def save_credentials(creds: CredentialsCreate):
    """Save LinkedIn credentials (will be encrypted)"""
    existing = await db.credentials.find_one({}, {"_id": 0})
    
    cred_obj = LinkedInCredentials(
        email=creds.email,
        password=creds.password,  # In production, encrypt this
        encrypted=False
    )
    doc = serialize_doc(cred_obj.model_dump())
    
    if existing:
        await db.credentials.update_one({"id": existing["id"]}, {"$set": doc})
    else:
        await db.credentials.insert_one(doc)
    
    await log_activity("credentials", "LinkedIn credentials saved", status="success")
    return {"message": "Credentials saved successfully", "email": creds.email}

@api_router.get("/credentials", response_model=dict)
async def get_credentials():
    """Get credentials status (not the actual password)"""
    creds = await db.credentials.find_one({}, {"_id": 0})
    if creds:
        return {
            "configured": True,
            "email": creds.get("email", ""),
            "updated_at": creds.get("updated_at")
        }
    return {"configured": False, "email": "", "updated_at": None}

@api_router.delete("/credentials")
async def delete_credentials():
    """Delete stored credentials"""
    await db.credentials.delete_many({})
    await log_activity("credentials", "LinkedIn credentials deleted", status="success")
    return {"message": "Credentials deleted"}

# ============== SEARCH CRITERIA ENDPOINTS ==============
@api_router.post("/search-criteria", response_model=SearchCriteria)
async def create_search_criteria(criteria: SearchCriteriaCreate):
    """Create new search criteria"""
    criteria_obj = SearchCriteria(**criteria.model_dump())
    doc = serialize_doc(criteria_obj.model_dump())
    await db.search_criteria.insert_one(doc)
    await log_activity("search", f"Created search criteria: {criteria.name}", status="success")
    return criteria_obj

@api_router.get("/search-criteria", response_model=List[SearchCriteria])
async def get_search_criteria():
    """Get all search criteria"""
    criteria_list = await db.search_criteria.find({}, {"_id": 0}).to_list(100)
    return [SearchCriteria(**deserialize_doc(c)) for c in criteria_list]

@api_router.put("/search-criteria/{criteria_id}", response_model=SearchCriteria)
async def update_search_criteria(criteria_id: str, criteria: SearchCriteriaCreate):
    """Update search criteria"""
    existing = await db.search_criteria.find_one({"id": criteria_id}, {"_id": 0})
    if not existing:
        raise HTTPException(status_code=404, detail="Search criteria not found")
    
    update_data = criteria.model_dump()
    await db.search_criteria.update_one({"id": criteria_id}, {"$set": update_data})
    
    updated = await db.search_criteria.find_one({"id": criteria_id}, {"_id": 0})
    return SearchCriteria(**deserialize_doc(updated))

@api_router.delete("/search-criteria/{criteria_id}")
async def delete_search_criteria(criteria_id: str):
    """Delete search criteria"""
    result = await db.search_criteria.delete_one({"id": criteria_id})
    if result.deleted_count == 0:
        raise HTTPException(status_code=404, detail="Search criteria not found")
    return {"message": "Search criteria deleted"}

# ============== MESSAGE TEMPLATES ENDPOINTS ==============
@api_router.post("/templates", response_model=MessageTemplate)
async def create_template(template: MessageTemplateCreate):
    """Create new message template"""
    template_obj = MessageTemplate(**template.model_dump())
    doc = serialize_doc(template_obj.model_dump())
    await db.templates.insert_one(doc)
    await log_activity("template", f"Created template: {template.name}", status="success")
    return template_obj

@api_router.get("/templates", response_model=List[MessageTemplate])
async def get_templates(template_type: Optional[str] = None):
    """Get all message templates"""
    query = {} if not template_type else {"template_type": template_type}
    templates = await db.templates.find(query, {"_id": 0}).to_list(100)
    return [MessageTemplate(**deserialize_doc(t)) for t in templates]

@api_router.put("/templates/{template_id}", response_model=MessageTemplate)
async def update_template(template_id: str, template: MessageTemplateCreate):
    """Update message template"""
    existing = await db.templates.find_one({"id": template_id}, {"_id": 0})
    if not existing:
        raise HTTPException(status_code=404, detail="Template not found")
    
    update_data = template.model_dump()
    await db.templates.update_one({"id": template_id}, {"$set": update_data})
    
    updated = await db.templates.find_one({"id": template_id}, {"_id": 0})
    return MessageTemplate(**deserialize_doc(updated))

@api_router.delete("/templates/{template_id}")
async def delete_template(template_id: str):
    """Delete message template"""
    result = await db.templates.delete_one({"id": template_id})
    if result.deleted_count == 0:
        raise HTTPException(status_code=404, detail="Template not found")
    return {"message": "Template deleted"}

# ============== RATE LIMIT CONFIG ENDPOINTS ==============
@api_router.get("/rate-limits", response_model=RateLimitConfig)
async def get_rate_limits():
    """Get rate limit configuration"""
    config = await db.rate_limits.find_one({}, {"_id": 0})
    if config:
        return RateLimitConfig(**deserialize_doc(config))
    
    # Return default config
    default_config = RateLimitConfig()
    doc = serialize_doc(default_config.model_dump())
    await db.rate_limits.insert_one(doc)
    return default_config

@api_router.put("/rate-limits", response_model=RateLimitConfig)
async def update_rate_limits(config: RateLimitConfigUpdate):
    """Update rate limit configuration"""
    existing = await db.rate_limits.find_one({}, {"_id": 0})
    
    update_data = {k: v for k, v in config.model_dump().items() if v is not None}
    update_data["updated_at"] = datetime.now(timezone.utc).isoformat()
    
    if existing:
        await db.rate_limits.update_one({"id": existing["id"]}, {"$set": update_data})
        updated = await db.rate_limits.find_one({"id": existing["id"]}, {"_id": 0})
    else:
        new_config = RateLimitConfig(**update_data)
        doc = serialize_doc(new_config.model_dump())
        await db.rate_limits.insert_one(doc)
        updated = doc
    
    await log_activity("config", "Rate limits updated", details=update_data, status="success")
    return RateLimitConfig(**deserialize_doc(updated))

# ============== STEALTH CONFIG ENDPOINTS ==============
@api_router.get("/stealth-config", response_model=StealthConfig)
async def get_stealth_config():
    """Get stealth/anti-bot configuration"""
    config = await db.stealth_config.find_one({}, {"_id": 0})
    if config:
        return StealthConfig(**deserialize_doc(config))
    
    # Return default config
    default_config = StealthConfig()
    doc = serialize_doc(default_config.model_dump())
    await db.stealth_config.insert_one(doc)
    return default_config

@api_router.put("/stealth-config", response_model=StealthConfig)
async def update_stealth_config(config: StealthConfigUpdate):
    """Update stealth/anti-bot configuration"""
    existing = await db.stealth_config.find_one({}, {"_id": 0})
    
    update_data = {k: v for k, v in config.model_dump().items() if v is not None}
    update_data["updated_at"] = datetime.now(timezone.utc).isoformat()
    
    if existing:
        await db.stealth_config.update_one({"id": existing["id"]}, {"$set": update_data})
        updated = await db.stealth_config.find_one({"id": existing["id"]}, {"_id": 0})
    else:
        new_config = StealthConfig(**update_data)
        doc = serialize_doc(new_config.model_dump())
        await db.stealth_config.insert_one(doc)
        updated = doc
    
    await log_activity("config", "Stealth configuration updated", status="success")
    return StealthConfig(**deserialize_doc(updated))

# ============== CONNECTIONS ENDPOINTS ==============
@api_router.get("/connections", response_model=List[Connection])
async def get_connections(status: Optional[str] = None, limit: int = 100, skip: int = 0):
    """Get all connections with optional status filter"""
    query = {} if not status else {"status": status}
    connections = await db.connections.find(query, {"_id": 0}).skip(skip).limit(limit).to_list(limit)
    return [Connection(**deserialize_doc(c)) for c in connections]

@api_router.get("/connections/{connection_id}", response_model=Connection)
async def get_connection(connection_id: str):
    """Get a specific connection"""
    conn = await db.connections.find_one({"id": connection_id}, {"_id": 0})
    if not conn:
        raise HTTPException(status_code=404, detail="Connection not found")
    return Connection(**deserialize_doc(conn))

@api_router.post("/connections", response_model=Connection)
async def create_connection(connection: dict):
    """Create a new connection record"""
    conn_obj = Connection(**connection)
    doc = serialize_doc(conn_obj.model_dump())
    await db.connections.insert_one(doc)
    return conn_obj

@api_router.put("/connections/{connection_id}/status")
async def update_connection_status(connection_id: str, status: str):
    """Update connection status"""
    update_data = {"status": status}
    if status == "accepted":
        update_data["accepted_at"] = datetime.now(timezone.utc).isoformat()
    
    result = await db.connections.update_one({"id": connection_id}, {"$set": update_data})
    if result.modified_count == 0:
        raise HTTPException(status_code=404, detail="Connection not found")
    return {"message": "Status updated"}

# ============== MESSAGES ENDPOINTS ==============
@api_router.get("/messages", response_model=List[Message])
async def get_messages(connection_id: Optional[str] = None, limit: int = 100):
    """Get all messages with optional connection filter"""
    query = {} if not connection_id else {"connection_id": connection_id}
    messages = await db.messages.find(query, {"_id": 0}).sort("sent_at", -1).to_list(limit)
    return [Message(**deserialize_doc(m)) for m in messages]

@api_router.post("/messages", response_model=Message)
async def create_message(message: dict):
    """Create a new message record"""
    msg_obj = Message(**message)
    doc = serialize_doc(msg_obj.model_dump())
    await db.messages.insert_one(doc)
    return msg_obj

# ============== ACTIVITY LOGS ENDPOINTS ==============
async def log_activity(action_type: str, description: str, details: dict = None, status: str = "success"):
    """Helper function to log activity"""
    log = ActivityLog(
        action_type=action_type,
        description=description,
        details=details,
        status=status
    )
    doc = serialize_doc(log.model_dump())
    await db.activity_logs.insert_one(doc)

@api_router.get("/activity-logs", response_model=List[ActivityLog])
async def get_activity_logs(action_type: Optional[str] = None, limit: int = 50):
    """Get activity logs"""
    query = {} if not action_type else {"action_type": action_type}
    logs = await db.activity_logs.find(query, {"_id": 0}).sort("timestamp", -1).to_list(limit)
    return [ActivityLog(**deserialize_doc(l)) for l in logs]

@api_router.delete("/activity-logs")
async def clear_activity_logs():
    """Clear all activity logs"""
    await db.activity_logs.delete_many({})
    return {"message": "Activity logs cleared"}

# ============== AUTOMATION STATE ENDPOINTS ==============
@api_router.get("/automation/state", response_model=AutomationState)
async def get_automation_state():
    """Get current automation state"""
    state = await db.automation_state.find_one({"id": "automation_state"}, {"_id": 0})
    if state:
        return AutomationState(**deserialize_doc(state))
    
    # Return default state
    default_state = AutomationState()
    doc = serialize_doc(default_state.model_dump())
    await db.automation_state.insert_one(doc)
    return default_state

@api_router.post("/automation/start")
async def start_automation(background_tasks: BackgroundTasks):
    """Start the automation process"""
    # Check if credentials are configured
    creds = await db.credentials.find_one({}, {"_id": 0})
    if not creds:
        raise HTTPException(status_code=400, detail="LinkedIn credentials not configured")
    
    # Update state
    await db.automation_state.update_one(
        {"id": "automation_state"},
        {"$set": {
            "status": AutomationStatus.RUNNING,
            "started_at": datetime.now(timezone.utc).isoformat(),
            "error_message": None
        }},
        upsert=True
    )
    
    await log_activity("automation", "Automation started", status="success")
    
    # In a real implementation, this would trigger the Go automation engine
    # For now, we'll simulate by updating state
    return {"message": "Automation started", "status": "running"}

@api_router.post("/automation/stop")
async def stop_automation():
    """Stop the automation process"""
    await db.automation_state.update_one(
        {"id": "automation_state"},
        {"$set": {
            "status": AutomationStatus.IDLE,
            "current_task": None
        }},
        upsert=True
    )
    
    await log_activity("automation", "Automation stopped", status="success")
    return {"message": "Automation stopped", "status": "idle"}

@api_router.post("/automation/pause")
async def pause_automation():
    """Pause the automation process"""
    await db.automation_state.update_one(
        {"id": "automation_state"},
        {"$set": {"status": AutomationStatus.PAUSED}},
        upsert=True
    )
    
    await log_activity("automation", "Automation paused", status="success")
    return {"message": "Automation paused", "status": "paused"}

@api_router.post("/automation/resume")
async def resume_automation():
    """Resume the automation process"""
    await db.automation_state.update_one(
        {"id": "automation_state"},
        {"$set": {"status": AutomationStatus.RUNNING}},
        upsert=True
    )
    
    await log_activity("automation", "Automation resumed", status="success")
    return {"message": "Automation resumed", "status": "running"}

# ============== DASHBOARD STATS ENDPOINT ==============
@api_router.get("/dashboard/stats", response_model=DashboardStats)
async def get_dashboard_stats():
    """Get dashboard statistics"""
    today = datetime.now(timezone.utc).replace(hour=0, minute=0, second=0, microsecond=0)
    
    # Count connections
    total_connections = await db.connections.count_documents({})
    pending_connections = await db.connections.count_documents({"status": "pending"})
    accepted_connections = await db.connections.count_documents({"status": "accepted"})
    
    # Count messages
    total_messages = await db.messages.count_documents({})
    messages_today = await db.messages.count_documents({
        "sent_at": {"$gte": today.isoformat()}
    })
    
    # Connections today
    connections_today = await db.connections.count_documents({
        "created_at": {"$gte": today.isoformat()}
    })
    
    # Get automation status
    state = await db.automation_state.find_one({"id": "automation_state"}, {"_id": 0})
    automation_status = AutomationStatus(state["status"]) if state else AutomationStatus.IDLE
    
    # Calculate success rate
    success_rate = 0.0
    if total_connections > 0:
        success_rate = (accepted_connections / total_connections) * 100
    
    return DashboardStats(
        total_connections=total_connections,
        pending_connections=pending_connections,
        accepted_connections=accepted_connections,
        total_messages=total_messages,
        messages_today=messages_today,
        connections_today=connections_today,
        automation_status=automation_status,
        success_rate=round(success_rate, 1)
    )

# ============== EXPORT CONFIG FOR GO ENGINE ==============
@api_router.get("/export-config")
async def export_config():
    """Export configuration for Go automation engine"""
    creds = await db.credentials.find_one({}, {"_id": 0})
    search_criteria = await db.search_criteria.find({"is_active": True}, {"_id": 0}).to_list(100)
    templates = await db.templates.find({"is_active": True}, {"_id": 0}).to_list(100)
    rate_limits = await db.rate_limits.find_one({}, {"_id": 0})
    stealth_config = await db.stealth_config.find_one({}, {"_id": 0})
    
    config = {
        "credentials": {
            "email": creds.get("email") if creds else "",
            "password": creds.get("password") if creds else ""
        } if creds else None,
        "search_criteria": search_criteria,
        "templates": {
            "connection": [t for t in templates if t.get("template_type") == "connection"],
            "follow_up": [t for t in templates if t.get("template_type") == "follow_up"]
        },
        "rate_limits": rate_limits or RateLimitConfig().model_dump(),
        "stealth": stealth_config or StealthConfig().model_dump()
    }
    
    return config

# ============== ROOT ENDPOINT ==============
@api_router.get("/")
async def root():
    return {"message": "LinkedIn Automation Dashboard API", "version": "1.0.0"}

# Include the router in the main app
app.include_router(api_router)

app.add_middleware(
    CORSMiddleware,
    allow_credentials=True,
    allow_origins=os.environ.get('CORS_ORIGINS', '*').split(','),
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.on_event("shutdown")
async def shutdown_db_client():
    client.close()
