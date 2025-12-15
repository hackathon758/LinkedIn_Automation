import { useState, useEffect, useCallback } from "react";
import "@/App.css";
import axios from "axios";
import { Toaster, toast } from "sonner";
import {
  Play,
  Pause,
  Square,
  Settings,
  Users,
  MessageSquare,
  Search,
  Shield,
  Activity,
  BarChart3,
  FileText,
  Key,
  Plus,
  Trash2,
  Edit,
  Check,
  X,
  RefreshCw,
  Clock,
  TrendingUp,
  UserPlus,
  Send,
  Eye,
  EyeOff,
  ChevronRight,
  AlertCircle,
  CheckCircle,
  Info,
  Zap,
  MousePointer,
  Timer,
  Fingerprint,
  Scroll,
  Keyboard,
  Calendar,
  Gauge,
  Globe,
  Wifi
} from "lucide-react";

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;
const API = `${BACKEND_URL}/api`;

// API Helper
const api = {
  get: (url) => axios.get(`${API}${url}`),
  post: (url, data) => axios.post(`${API}${url}`, data),
  put: (url, data) => axios.put(`${API}${url}`, data),
  delete: (url) => axios.delete(`${API}${url}`),
};

// ============== COMPONENTS ==============

// Sidebar Navigation
const Sidebar = ({ activeTab, setActiveTab }) => {
  const menuItems = [
    { id: "dashboard", label: "Dashboard", icon: BarChart3 },
    { id: "credentials", label: "Credentials", icon: Key },
    { id: "search", label: "Search Criteria", icon: Search },
    { id: "templates", label: "Templates", icon: FileText },
    { id: "connections", label: "Connections", icon: Users },
    { id: "messages", label: "Messages", icon: MessageSquare },
    { id: "stealth", label: "Anti-Bot Settings", icon: Shield },
    { id: "rate-limits", label: "Rate Limits", icon: Gauge },
    { id: "activity", label: "Activity Logs", icon: Activity },
  ];

  return (
    <div className="w-64 bg-slate-900 min-h-screen p-4 flex flex-col" data-testid="sidebar">
      <div className="mb-8">
        <h1 className="text-xl font-bold text-white flex items-center gap-2">
          <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <Zap className="w-5 h-5 text-white" />
          </div>
          LinkedIn Auto
        </h1>
        <p className="text-slate-400 text-sm mt-1">Automation Dashboard</p>
      </div>

      <nav className="flex-1 space-y-1">
        {menuItems.map((item) => {
          const Icon = item.icon;
          return (
            <button
              key={item.id}
              data-testid={`nav-${item.id}`}
              onClick={() => setActiveTab(item.id)}
              className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-all ${
                activeTab === item.id
                  ? "bg-blue-600 text-white"
                  : "text-slate-300 hover:bg-slate-800 hover:text-white"
              }`}
            >
              <Icon className="w-5 h-5" />
              <span className="font-medium">{item.label}</span>
            </button>
          );
        })}
      </nav>

      <div className="mt-auto pt-4 border-t border-slate-700">
        <div className="text-slate-400 text-xs">
          <p>Version 1.0.0</p>
          <p className="mt-1">Go + Rod Engine</p>
        </div>
      </div>
    </div>
  );
};

// Stats Card
const StatsCard = ({ title, value, icon: Icon, trend, color = "blue" }) => {
  const colorClasses = {
    blue: "bg-blue-500/10 text-blue-500",
    green: "bg-green-500/10 text-green-500",
    yellow: "bg-yellow-500/10 text-yellow-500",
    purple: "bg-purple-500/10 text-purple-500",
    red: "bg-red-500/10 text-red-500",
  };

  return (
    <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200" data-testid={`stats-${title.toLowerCase().replace(/\s/g, '-')}`}>
      <div className="flex items-start justify-between">
        <div>
          <p className="text-slate-500 text-sm font-medium">{title}</p>
          <p className="text-2xl font-bold text-slate-900 mt-1">{value}</p>
          {trend && (
            <p className="text-green-500 text-sm mt-1 flex items-center gap-1">
              <TrendingUp className="w-3 h-3" /> {trend}
            </p>
          )}
        </div>
        <div className={`p-3 rounded-xl ${colorClasses[color]}`}>
          <Icon className="w-6 h-6" />
        </div>
      </div>
    </div>
  );
};

// Automation Control Panel
const AutomationControls = ({ status, onStart, onStop, onPause, onResume }) => {
  const statusColors = {
    idle: "bg-slate-100 text-slate-600",
    running: "bg-green-100 text-green-700",
    paused: "bg-yellow-100 text-yellow-700",
    error: "bg-red-100 text-red-700",
  };

  const statusLabels = {
    idle: "Idle",
    running: "Running",
    paused: "Paused",
    error: "Error",
  };

  return (
    <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200" data-testid="automation-controls">
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-semibold text-slate-900">Automation Control</h3>
        <span className={`px-3 py-1 rounded-full text-sm font-medium ${statusColors[status]}`}>
          {statusLabels[status]}
        </span>
      </div>
      <div className="flex gap-2">
        {status === "idle" && (
          <button
            data-testid="btn-start"
            onClick={onStart}
            className="flex-1 bg-green-600 hover:bg-green-700 text-white py-2.5 px-4 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors"
          >
            <Play className="w-4 h-4" /> Start
          </button>
        )}
        {status === "running" && (
          <>
            <button
              data-testid="btn-pause"
              onClick={onPause}
              className="flex-1 bg-yellow-500 hover:bg-yellow-600 text-white py-2.5 px-4 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors"
            >
              <Pause className="w-4 h-4" /> Pause
            </button>
            <button
              data-testid="btn-stop"
              onClick={onStop}
              className="flex-1 bg-red-500 hover:bg-red-600 text-white py-2.5 px-4 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors"
            >
              <Square className="w-4 h-4" /> Stop
            </button>
          </>
        )}
        {status === "paused" && (
          <>
            <button
              data-testid="btn-resume"
              onClick={onResume}
              className="flex-1 bg-green-600 hover:bg-green-700 text-white py-2.5 px-4 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors"
            >
              <Play className="w-4 h-4" /> Resume
            </button>
            <button
              data-testid="btn-stop-paused"
              onClick={onStop}
              className="flex-1 bg-red-500 hover:bg-red-600 text-white py-2.5 px-4 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors"
            >
              <Square className="w-4 h-4" /> Stop
            </button>
          </>
        )}
        {status === "error" && (
          <button
            data-testid="btn-restart"
            onClick={onStart}
            className="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-2.5 px-4 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors"
          >
            <RefreshCw className="w-4 h-4" /> Restart
          </button>
        )}
      </div>
    </div>
  );
};

// Dashboard Page
const DashboardPage = () => {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);

  const fetchStats = useCallback(async () => {
    try {
      const response = await api.get("/dashboard/stats");
      setStats(response.data);
    } catch (error) {
      console.error("Failed to fetch stats", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 10000);
    return () => clearInterval(interval);
  }, [fetchStats]);

  const handleStart = async () => {
    try {
      await api.post("/automation/start");
      toast.success("Automation started");
      fetchStats();
    } catch (error) {
      toast.error(error.response?.data?.detail || "Failed to start automation");
    }
  };

  const handleStop = async () => {
    try {
      await api.post("/automation/stop");
      toast.success("Automation stopped");
      fetchStats();
    } catch (error) {
      toast.error("Failed to stop automation");
    }
  };

  const handlePause = async () => {
    try {
      await api.post("/automation/pause");
      toast.success("Automation paused");
      fetchStats();
    } catch (error) {
      toast.error("Failed to pause automation");
    }
  };

  const handleResume = async () => {
    try {
      await api.post("/automation/resume");
      toast.success("Automation resumed");
      fetchStats();
    } catch (error) {
      toast.error("Failed to resume automation");
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="space-y-6" data-testid="dashboard-page">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Dashboard</h2>
          <p className="text-slate-500">Monitor your LinkedIn automation</p>
        </div>
        <button
          onClick={fetchStats}
          className="p-2 hover:bg-slate-100 rounded-lg transition-colors"
          data-testid="btn-refresh"
        >
          <RefreshCw className="w-5 h-5 text-slate-600" />
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="Total Connections"
          value={stats?.total_connections || 0}
          icon={Users}
          color="blue"
        />
        <StatsCard
          title="Accepted"
          value={stats?.accepted_connections || 0}
          icon={UserPlus}
          color="green"
        />
        <StatsCard
          title="Messages Sent"
          value={stats?.total_messages || 0}
          icon={Send}
          color="purple"
        />
        <StatsCard
          title="Success Rate"
          value={`${stats?.success_rate || 0}%`}
          icon={TrendingUp}
          color="yellow"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <AutomationControls
          status={stats?.automation_status || "idle"}
          onStart={handleStart}
          onStop={handleStop}
          onPause={handlePause}
          onResume={handleResume}
        />

        <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
          <h3 className="font-semibold text-slate-900 mb-4">Today's Activity</h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-slate-600">Connections Sent</span>
              <span className="font-semibold text-slate-900">{stats?.connections_today || 0}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-slate-600">Messages Sent</span>
              <span className="font-semibold text-slate-900">{stats?.messages_today || 0}</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-slate-600">Pending Connections</span>
              <span className="font-semibold text-slate-900">{stats?.pending_connections || 0}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Credentials Page
const CredentialsPage = () => {
  const [credentials, setCredentials] = useState({ configured: false, email: "" });
  const [formData, setFormData] = useState({ email: "", password: "" });
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetchCredentials();
  }, []);

  const fetchCredentials = async () => {
    try {
      const response = await api.get("/credentials");
      setCredentials(response.data);
      if (response.data.email) {
        setFormData((prev) => ({ ...prev, email: response.data.email }));
      }
    } catch (error) {
      console.error("Failed to fetch credentials", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    if (!formData.email || !formData.password) {
      toast.error("Please fill in all fields");
      return;
    }

    setSaving(true);
    try {
      await api.post("/credentials", formData);
      toast.success("Credentials saved successfully");
      setFormData((prev) => ({ ...prev, password: "" }));
      fetchCredentials();
    } catch (error) {
      toast.error("Failed to save credentials");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm("Are you sure you want to delete your credentials?")) return;

    try {
      await api.delete("/credentials");
      toast.success("Credentials deleted");
      setCredentials({ configured: false, email: "" });
      setFormData({ email: "", password: "" });
    } catch (error) {
      toast.error("Failed to delete credentials");
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="max-w-2xl" data-testid="credentials-page">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-slate-900">LinkedIn Credentials</h2>
        <p className="text-slate-500">Configure your LinkedIn login credentials</p>
      </div>

      <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200">
        {credentials.configured && (
          <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg flex items-center gap-3">
            <CheckCircle className="w-5 h-5 text-green-600" />
            <div>
              <p className="font-medium text-green-800">Credentials Configured</p>
              <p className="text-sm text-green-600">Email: {credentials.email}</p>
            </div>
          </div>
        )}

        <form onSubmit={handleSave} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              LinkedIn Email
            </label>
            <input
              type="email"
              data-testid="input-email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
              placeholder="your-email@example.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              LinkedIn Password
            </label>
            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                data-testid="input-password"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none pr-10"
                placeholder="••••••••"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
              >
                {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="submit"
              data-testid="btn-save-credentials"
              disabled={saving}
              className="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-2.5 px-4 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors disabled:opacity-50"
            >
              {saving ? <RefreshCw className="w-4 h-4 animate-spin" /> : <Check className="w-4 h-4" />}
              {credentials.configured ? "Update Credentials" : "Save Credentials"}
            </button>
            {credentials.configured && (
              <button
                type="button"
                onClick={handleDelete}
                data-testid="btn-delete-credentials"
                className="bg-red-50 hover:bg-red-100 text-red-600 py-2.5 px-4 rounded-lg font-medium flex items-center justify-center gap-2 transition-colors"
              >
                <Trash2 className="w-4 h-4" /> Delete
              </button>
            )}
          </div>
        </form>

        <div className="mt-6 p-4 bg-amber-50 border border-amber-200 rounded-lg">
          <div className="flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-amber-600 mt-0.5" />
            <div>
              <p className="font-medium text-amber-800">Security Notice</p>
              <p className="text-sm text-amber-700 mt-1">
                Your credentials are stored securely and only used for LinkedIn automation.
                We recommend using a dedicated LinkedIn account for automation purposes.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Search Criteria Page
const SearchCriteriaPage = () => {
  const [criteria, setCriteria] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [formData, setFormData] = useState({
    name: "",
    job_titles: "",
    companies: "",
    locations: "",
    keywords: "",
  });

  useEffect(() => {
    fetchCriteria();
  }, []);

  const fetchCriteria = async () => {
    try {
      const response = await api.get("/search-criteria");
      setCriteria(response.data);
    } catch (error) {
      console.error("Failed to fetch criteria", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    if (!formData.name) {
      toast.error("Please enter a name");
      return;
    }

    const payload = {
      name: formData.name,
      job_titles: formData.job_titles.split(",").map((s) => s.trim()).filter(Boolean),
      companies: formData.companies.split(",").map((s) => s.trim()).filter(Boolean),
      locations: formData.locations.split(",").map((s) => s.trim()).filter(Boolean),
      keywords: formData.keywords.split(",").map((s) => s.trim()).filter(Boolean),
    };

    try {
      if (editingId) {
        await api.put(`/search-criteria/${editingId}`, payload);
        toast.success("Search criteria updated");
      } else {
        await api.post("/search-criteria", payload);
        toast.success("Search criteria created");
      }
      resetForm();
      fetchCriteria();
    } catch (error) {
      toast.error("Failed to save search criteria");
    }
  };

  const handleEdit = (item) => {
    setFormData({
      name: item.name,
      job_titles: item.job_titles.join(", "),
      companies: item.companies.join(", "),
      locations: item.locations.join(", "),
      keywords: item.keywords.join(", "),
    });
    setEditingId(item.id);
    setShowForm(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm("Delete this search criteria?")) return;
    try {
      await api.delete(`/search-criteria/${id}`);
      toast.success("Search criteria deleted");
      fetchCriteria();
    } catch (error) {
      toast.error("Failed to delete");
    }
  };

  const resetForm = () => {
    setFormData({ name: "", job_titles: "", companies: "", locations: "", keywords: "" });
    setEditingId(null);
    setShowForm(false);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div data-testid="search-criteria-page">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Search Criteria</h2>
          <p className="text-slate-500">Define who to target for connections</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          data-testid="btn-add-criteria"
          className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg font-medium flex items-center gap-2 transition-colors"
        >
          <Plus className="w-4 h-4" /> Add Criteria
        </button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200 mb-6">
          <h3 className="font-semibold text-slate-900 mb-4">
            {editingId ? "Edit Search Criteria" : "New Search Criteria"}
          </h3>
          <form onSubmit={handleSave} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">Name</label>
              <input
                type="text"
                data-testid="input-criteria-name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                placeholder="e.g., Tech Recruiters"
              />
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Job Titles (comma-separated)
                </label>
                <input
                  type="text"
                  data-testid="input-job-titles"
                  value={formData.job_titles}
                  onChange={(e) => setFormData({ ...formData, job_titles: e.target.value })}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                  placeholder="Recruiter, HR Manager"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Companies (comma-separated)
                </label>
                <input
                  type="text"
                  data-testid="input-companies"
                  value={formData.companies}
                  onChange={(e) => setFormData({ ...formData, companies: e.target.value })}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                  placeholder="Google, Microsoft"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Locations (comma-separated)
                </label>
                <input
                  type="text"
                  data-testid="input-locations"
                  value={formData.locations}
                  onChange={(e) => setFormData({ ...formData, locations: e.target.value })}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                  placeholder="San Francisco, New York"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Keywords (comma-separated)
                </label>
                <input
                  type="text"
                  data-testid="input-keywords"
                  value={formData.keywords}
                  onChange={(e) => setFormData({ ...formData, keywords: e.target.value })}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                  placeholder="hiring, tech"
                />
              </div>
            </div>
            <div className="flex gap-3 pt-2">
              <button
                type="submit"
                data-testid="btn-save-criteria"
                className="bg-blue-600 hover:bg-blue-700 text-white py-2.5 px-4 rounded-lg font-medium flex items-center gap-2 transition-colors"
              >
                <Check className="w-4 h-4" /> {editingId ? "Update" : "Save"}
              </button>
              <button
                type="button"
                onClick={resetForm}
                className="bg-slate-100 hover:bg-slate-200 text-slate-700 py-2.5 px-4 rounded-lg font-medium flex items-center gap-2 transition-colors"
              >
                <X className="w-4 h-4" /> Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="space-y-4">
        {criteria.length === 0 ? (
          <div className="bg-white rounded-xl p-8 shadow-sm border border-slate-200 text-center">
            <Search className="w-12 h-12 text-slate-300 mx-auto mb-3" />
            <p className="text-slate-500">No search criteria defined yet</p>
            <p className="text-sm text-slate-400 mt-1">Create your first search criteria to start targeting users</p>
          </div>
        ) : (
          criteria.map((item) => (
            <div key={item.id} className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
              <div className="flex items-start justify-between">
                <div>
                  <h4 className="font-semibold text-slate-900">{item.name}</h4>
                  <div className="mt-3 space-y-2">
                    {item.job_titles.length > 0 && (
                      <div className="flex items-center gap-2 text-sm">
                        <span className="text-slate-500">Job Titles:</span>
                        <span className="text-slate-700">{item.job_titles.join(", ")}</span>
                      </div>
                    )}
                    {item.companies.length > 0 && (
                      <div className="flex items-center gap-2 text-sm">
                        <span className="text-slate-500">Companies:</span>
                        <span className="text-slate-700">{item.companies.join(", ")}</span>
                      </div>
                    )}
                    {item.locations.length > 0 && (
                      <div className="flex items-center gap-2 text-sm">
                        <span className="text-slate-500">Locations:</span>
                        <span className="text-slate-700">{item.locations.join(", ")}</span>
                      </div>
                    )}
                    {item.keywords.length > 0 && (
                      <div className="flex items-center gap-2 text-sm">
                        <span className="text-slate-500">Keywords:</span>
                        <span className="text-slate-700">{item.keywords.join(", ")}</span>
                      </div>
                    )}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => handleEdit(item)}
                    data-testid={`btn-edit-${item.id}`}
                    className="p-2 hover:bg-slate-100 rounded-lg transition-colors"
                  >
                    <Edit className="w-4 h-4 text-slate-500" />
                  </button>
                  <button
                    onClick={() => handleDelete(item.id)}
                    data-testid={`btn-delete-${item.id}`}
                    className="p-2 hover:bg-red-50 rounded-lg transition-colors"
                  >
                    <Trash2 className="w-4 h-4 text-red-500" />
                  </button>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

// Templates Page
const TemplatesPage = () => {
  const [templates, setTemplates] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [formData, setFormData] = useState({
    name: "",
    template_type: "connection",
    content: "",
    variables: "",
  });

  useEffect(() => {
    fetchTemplates();
  }, []);

  const fetchTemplates = async () => {
    try {
      const response = await api.get("/templates");
      setTemplates(response.data);
    } catch (error) {
      console.error("Failed to fetch templates", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    if (!formData.name || !formData.content) {
      toast.error("Please fill in required fields");
      return;
    }

    const payload = {
      name: formData.name,
      template_type: formData.template_type,
      content: formData.content,
      variables: formData.variables.split(",").map((s) => s.trim()).filter(Boolean),
    };

    try {
      if (editingId) {
        await api.put(`/templates/${editingId}`, payload);
        toast.success("Template updated");
      } else {
        await api.post("/templates", payload);
        toast.success("Template created");
      }
      resetForm();
      fetchTemplates();
    } catch (error) {
      toast.error("Failed to save template");
    }
  };

  const handleEdit = (item) => {
    setFormData({
      name: item.name,
      template_type: item.template_type,
      content: item.content,
      variables: item.variables.join(", "),
    });
    setEditingId(item.id);
    setShowForm(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm("Delete this template?")) return;
    try {
      await api.delete(`/templates/${id}`);
      toast.success("Template deleted");
      fetchTemplates();
    } catch (error) {
      toast.error("Failed to delete");
    }
  };

  const resetForm = () => {
    setFormData({ name: "", template_type: "connection", content: "", variables: "" });
    setEditingId(null);
    setShowForm(false);
  };

  const connectionTemplates = templates.filter((t) => t.template_type === "connection");
  const followUpTemplates = templates.filter((t) => t.template_type === "follow_up");

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div data-testid="templates-page">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Message Templates</h2>
          <p className="text-slate-500">Create personalized message templates</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          data-testid="btn-add-template"
          className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg font-medium flex items-center gap-2 transition-colors"
        >
          <Plus className="w-4 h-4" /> Add Template
        </button>
      </div>

      {showForm && (
        <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200 mb-6">
          <h3 className="font-semibold text-slate-900 mb-4">
            {editingId ? "Edit Template" : "New Template"}
          </h3>
          <form onSubmit={handleSave} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Name</label>
                <input
                  type="text"
                  data-testid="input-template-name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                  placeholder="e.g., Professional Introduction"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Type</label>
                <select
                  data-testid="select-template-type"
                  value={formData.template_type}
                  onChange={(e) => setFormData({ ...formData, template_type: e.target.value })}
                  className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                >
                  <option value="connection">Connection Request</option>
                  <option value="follow_up">Follow-up Message</option>
                </select>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Message Content
              </label>
              <textarea
                data-testid="textarea-template-content"
                value={formData.content}
                onChange={(e) => setFormData({ ...formData, content: e.target.value })}
                className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none min-h-[120px]"
                placeholder="Hi {{firstName}}, I noticed you work at {{company}}..."
              />
              <p className="text-xs text-slate-500 mt-1">
                Use variables like {`{{firstName}}`}, {`{{lastName}}`}, {`{{company}}`}, {`{{jobTitle}}`}
              </p>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Variables (comma-separated)
              </label>
              <input
                type="text"
                data-testid="input-template-variables"
                value={formData.variables}
                onChange={(e) => setFormData({ ...formData, variables: e.target.value })}
                className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                placeholder="firstName, company, jobTitle"
              />
            </div>
            <div className="flex gap-3 pt-2">
              <button
                type="submit"
                data-testid="btn-save-template"
                className="bg-blue-600 hover:bg-blue-700 text-white py-2.5 px-4 rounded-lg font-medium flex items-center gap-2 transition-colors"
              >
                <Check className="w-4 h-4" /> {editingId ? "Update" : "Save"}
              </button>
              <button
                type="button"
                onClick={resetForm}
                className="bg-slate-100 hover:bg-slate-200 text-slate-700 py-2.5 px-4 rounded-lg font-medium flex items-center gap-2 transition-colors"
              >
                <X className="w-4 h-4" /> Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="space-y-6">
        <div>
          <h3 className="font-semibold text-slate-900 mb-3 flex items-center gap-2">
            <UserPlus className="w-5 h-5" /> Connection Request Templates
          </h3>
          {connectionTemplates.length === 0 ? (
            <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200 text-center">
              <p className="text-slate-500">No connection templates yet</p>
            </div>
          ) : (
            <div className="space-y-3">
              {connectionTemplates.map((item) => (
                <TemplateCard key={item.id} item={item} onEdit={handleEdit} onDelete={handleDelete} />
              ))}
            </div>
          )}
        </div>

        <div>
          <h3 className="font-semibold text-slate-900 mb-3 flex items-center gap-2">
            <Send className="w-5 h-5" /> Follow-up Message Templates
          </h3>
          {followUpTemplates.length === 0 ? (
            <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200 text-center">
              <p className="text-slate-500">No follow-up templates yet</p>
            </div>
          ) : (
            <div className="space-y-3">
              {followUpTemplates.map((item) => (
                <TemplateCard key={item.id} item={item} onEdit={handleEdit} onDelete={handleDelete} />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const TemplateCard = ({ item, onEdit, onDelete }) => (
  <div className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
    <div className="flex items-start justify-between">
      <div className="flex-1">
        <div className="flex items-center gap-2">
          <h4 className="font-semibold text-slate-900">{item.name}</h4>
          <span className="text-xs bg-slate-100 text-slate-600 px-2 py-0.5 rounded">
            Used {item.usage_count}x
          </span>
        </div>
        <p className="text-sm text-slate-600 mt-2 whitespace-pre-wrap">{item.content}</p>
        {item.variables.length > 0 && (
          <div className="flex flex-wrap gap-1 mt-3">
            {item.variables.map((v) => (
              <span key={v} className="text-xs bg-blue-50 text-blue-600 px-2 py-0.5 rounded">
                {`{{${v}}}`}
              </span>
            ))}
          </div>
        )}
      </div>
      <div className="flex items-center gap-2 ml-4">
        <button onClick={() => onEdit(item)} className="p-2 hover:bg-slate-100 rounded-lg transition-colors">
          <Edit className="w-4 h-4 text-slate-500" />
        </button>
        <button onClick={() => onDelete(item.id)} className="p-2 hover:bg-red-50 rounded-lg transition-colors">
          <Trash2 className="w-4 h-4 text-red-500" />
        </button>
      </div>
    </div>
  </div>
);

// Connections Page
const ConnectionsPage = () => {
  const [connections, setConnections] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState("");

  useEffect(() => {
    fetchConnections();
  }, [filter]);

  const fetchConnections = async () => {
    setLoading(true);
    try {
      const url = filter ? `/connections?status=${filter}` : "/connections";
      const response = await api.get(url);
      setConnections(response.data);
    } catch (error) {
      console.error("Failed to fetch connections", error);
    } finally {
      setLoading(false);
    }
  };

  const statusBadge = (status) => {
    const colors = {
      pending: "bg-yellow-100 text-yellow-700",
      accepted: "bg-green-100 text-green-700",
      declined: "bg-red-100 text-red-700",
      failed: "bg-slate-100 text-slate-600",
    };
    return (
      <span className={`px-2 py-0.5 rounded text-xs font-medium ${colors[status]}`}>
        {status}
      </span>
    );
  };

  return (
    <div data-testid="connections-page">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Connections</h2>
          <p className="text-slate-500">View and manage your connection requests</p>
        </div>
        <div className="flex items-center gap-3">
          <select
            data-testid="filter-status"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
          >
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="accepted">Accepted</option>
            <option value="declined">Declined</option>
            <option value="failed">Failed</option>
          </select>
          <button
            onClick={fetchConnections}
            className="p-2 hover:bg-slate-100 rounded-lg transition-colors"
          >
            <RefreshCw className="w-5 h-5 text-slate-600" />
          </button>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
        </div>
      ) : connections.length === 0 ? (
        <div className="bg-white rounded-xl p-8 shadow-sm border border-slate-200 text-center">
          <Users className="w-12 h-12 text-slate-300 mx-auto mb-3" />
          <p className="text-slate-500">No connections found</p>
          <p className="text-sm text-slate-400 mt-1">Start the automation to send connection requests</p>
        </div>
      ) : (
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
          <table className="w-full">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="text-left py-3 px-4 text-sm font-medium text-slate-600">Profile</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-slate-600">Company</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-slate-600">Status</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-slate-600">Date</th>
              </tr>
            </thead>
            <tbody>
              {connections.map((conn) => (
                <tr key={conn.id} className="border-b border-slate-100 hover:bg-slate-50">
                  <td className="py-3 px-4">
                    <div>
                      <p className="font-medium text-slate-900">
                        {conn.first_name} {conn.last_name}
                      </p>
                      <p className="text-sm text-slate-500">{conn.job_title}</p>
                    </div>
                  </td>
                  <td className="py-3 px-4 text-slate-600">{conn.company || "-"}</td>
                  <td className="py-3 px-4">{statusBadge(conn.status)}</td>
                  <td className="py-3 px-4 text-slate-500 text-sm">
                    {new Date(conn.created_at).toLocaleDateString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

// Messages Page
const MessagesPage = () => {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchMessages();
  }, []);

  const fetchMessages = async () => {
    try {
      const response = await api.get("/messages");
      setMessages(response.data);
    } catch (error) {
      console.error("Failed to fetch messages", error);
    } finally {
      setLoading(false);
    }
  };

  const statusBadge = (status) => {
    const colors = {
      sent: "bg-blue-100 text-blue-700",
      delivered: "bg-green-100 text-green-700",
      read: "bg-purple-100 text-purple-700",
      failed: "bg-red-100 text-red-700",
    };
    return (
      <span className={`px-2 py-0.5 rounded text-xs font-medium ${colors[status]}`}>
        {status}
      </span>
    );
  };

  return (
    <div data-testid="messages-page">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Messages</h2>
          <p className="text-slate-500">View sent follow-up messages</p>
        </div>
        <button
          onClick={fetchMessages}
          className="p-2 hover:bg-slate-100 rounded-lg transition-colors"
        >
          <RefreshCw className="w-5 h-5 text-slate-600" />
        </button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
        </div>
      ) : messages.length === 0 ? (
        <div className="bg-white rounded-xl p-8 shadow-sm border border-slate-200 text-center">
          <MessageSquare className="w-12 h-12 text-slate-300 mx-auto mb-3" />
          <p className="text-slate-500">No messages sent yet</p>
          <p className="text-sm text-slate-400 mt-1">Messages will appear here after follow-up automation</p>
        </div>
      ) : (
        <div className="space-y-4">
          {messages.map((msg) => (
            <div key={msg.id} className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
              <div className="flex items-start justify-between">
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    {statusBadge(msg.status)}
                    <span className="text-sm text-slate-500">
                      {new Date(msg.sent_at).toLocaleString()}
                    </span>
                  </div>
                  <p className="text-slate-700">{msg.content}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

// Stealth Config Page
const StealthConfigPage = () => {
  const [config, setConfig] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    try {
      const response = await api.get("/stealth-config");
      setConfig(response.data);
    } catch (error) {
      console.error("Failed to fetch stealth config", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await api.put("/stealth-config", config);
      toast.success("Anti-bot settings saved");
    } catch (error) {
      toast.error("Failed to save settings");
    } finally {
      setSaving(false);
    }
  };

  const updateConfig = (key, value) => {
    setConfig({ ...config, [key]: value });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    );
  }

  const sections = [
    {
      title: "Bézier Curve Mouse Movement",
      icon: MousePointer,
      description: "Simulate natural mouse movement patterns",
      mandatory: true,
      settings: [
        { key: "bezier_enabled", label: "Enable Bézier Curves", type: "toggle" },
        { key: "bezier_overshoot_probability", label: "Overshoot Probability", type: "slider", min: 0, max: 1, step: 0.05 },
      ],
    },
    {
      title: "Randomized Timing Patterns",
      icon: Timer,
      description: "Add realistic delays between actions",
      mandatory: true,
      settings: [
        { key: "typing_min_delay_ms", label: "Min Typing Delay (ms)", type: "number" },
        { key: "typing_max_delay_ms", label: "Max Typing Delay (ms)", type: "number" },
        { key: "typo_probability", label: "Typo Probability", type: "slider", min: 0, max: 0.2, step: 0.01 },
      ],
    },
    {
      title: "Browser Fingerprint Masking",
      icon: Fingerprint,
      description: "Evade fingerprint-based detection",
      mandatory: true,
      settings: [
        { key: "rotate_user_agent", label: "Rotate User Agent", type: "toggle" },
        { key: "randomize_viewport", label: "Randomize Viewport", type: "toggle" },
        { key: "disable_webdriver_flag", label: "Disable Webdriver Flag", type: "toggle" },
      ],
    },
    {
      title: "Random Scrolling Behavior",
      icon: Scroll,
      description: "Natural scrolling patterns",
      settings: [
        { key: "scroll_min_speed", label: "Min Scroll Speed (px)", type: "number" },
        { key: "scroll_max_speed", label: "Max Scroll Speed (px)", type: "number" },
        { key: "scroll_back_probability", label: "Scroll Back Probability", type: "slider", min: 0, max: 0.3, step: 0.05 },
      ],
    },
    {
      title: "Realistic Typing Simulation",
      icon: Keyboard,
      description: "Human-like typing patterns with typos",
      settings: [
        { key: "hover_before_click", label: "Hover Before Click", type: "toggle" },
        { key: "random_cursor_movement", label: "Random Cursor Movement", type: "toggle" },
      ],
    },
    {
      title: "Activity Scheduling",
      icon: Calendar,
      description: "Operate during realistic hours",
      settings: [
        { key: "respect_business_hours", label: "Respect Business Hours", type: "toggle" },
        { key: "include_break_patterns", label: "Include Break Patterns", type: "toggle" },
      ],
    },
    {
      title: "Rate Limiting & Throttling",
      icon: Gauge,
      description: "Control action frequency",
      settings: [
        { key: "enable_token_bucket", label: "Enable Token Bucket", type: "toggle" },
        { key: "cooldown_after_bulk", label: "Cooldown After Bulk", type: "toggle" },
      ],
    },
    {
      title: "Request Headers & Network",
      icon: Globe,
      description: "Vary HTTP headers and network behavior",
      settings: [
        { key: "randomize_headers", label: "Randomize Headers", type: "toggle" },
        { key: "simulate_network_latency", label: "Simulate Network Latency", type: "toggle" },
      ],
    },
  ];

  return (
    <div data-testid="stealth-config-page">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Anti-Bot Detection Settings</h2>
          <p className="text-slate-500">Configure stealth techniques to evade detection</p>
        </div>
        <button
          onClick={handleSave}
          disabled={saving}
          data-testid="btn-save-stealth"
          className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg font-medium flex items-center gap-2 transition-colors disabled:opacity-50"
        >
          {saving ? <RefreshCw className="w-4 h-4 animate-spin" /> : <Check className="w-4 h-4" />}
          Save Settings
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {sections.map((section) => {
          const Icon = section.icon;
          return (
            <div key={section.title} className="bg-white rounded-xl p-5 shadow-sm border border-slate-200">
              <div className="flex items-start gap-3 mb-4">
                <div className={`p-2 rounded-lg ${section.mandatory ? 'bg-red-50' : 'bg-blue-50'}`}>
                  <Icon className={`w-5 h-5 ${section.mandatory ? 'text-red-600' : 'text-blue-600'}`} />
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold text-slate-900">{section.title}</h3>
                    {section.mandatory && (
                      <span className="text-xs bg-red-100 text-red-600 px-2 py-0.5 rounded">Required</span>
                    )}
                  </div>
                  <p className="text-sm text-slate-500">{section.description}</p>
                </div>
              </div>
              <div className="space-y-4">
                {section.settings.map((setting) => (
                  <div key={setting.key} className="flex items-center justify-between">
                    <span className="text-sm text-slate-700">{setting.label}</span>
                    {setting.type === "toggle" && (
                      <button
                        onClick={() => updateConfig(setting.key, !config[setting.key])}
                        className={`w-11 h-6 rounded-full transition-colors ${
                          config[setting.key] ? "bg-blue-600" : "bg-slate-200"
                        }`}
                      >
                        <div
                          className={`w-5 h-5 rounded-full bg-white shadow transform transition-transform ${
                            config[setting.key] ? "translate-x-5" : "translate-x-0.5"
                          }`}
                        />
                      </button>
                    )}
                    {setting.type === "number" && (
                      <input
                        type="number"
                        value={config[setting.key]}
                        onChange={(e) => updateConfig(setting.key, parseInt(e.target.value) || 0)}
                        className="w-24 px-3 py-1.5 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                      />
                    )}
                    {setting.type === "slider" && (
                      <div className="flex items-center gap-2">
                        <input
                          type="range"
                          min={setting.min}
                          max={setting.max}
                          step={setting.step}
                          value={config[setting.key]}
                          onChange={(e) => updateConfig(setting.key, parseFloat(e.target.value))}
                          className="w-24"
                        />
                        <span className="text-sm text-slate-600 w-12">
                          {(config[setting.key] * 100).toFixed(0)}%
                        </span>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

// Rate Limits Page
const RateLimitsPage = () => {
  const [config, setConfig] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    try {
      const response = await api.get("/rate-limits");
      setConfig(response.data);
    } catch (error) {
      console.error("Failed to fetch rate limits", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await api.put("/rate-limits", config);
      toast.success("Rate limits saved");
    } catch (error) {
      toast.error("Failed to save rate limits");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="max-w-2xl" data-testid="rate-limits-page">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Rate Limits</h2>
          <p className="text-slate-500">Configure daily limits and timing</p>
        </div>
        <button
          onClick={handleSave}
          disabled={saving}
          data-testid="btn-save-rate-limits"
          className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg font-medium flex items-center gap-2 transition-colors disabled:opacity-50"
        >
          {saving ? <RefreshCw className="w-4 h-4 animate-spin" /> : <Check className="w-4 h-4" />}
          Save
        </button>
      </div>

      <div className="bg-white rounded-xl p-6 shadow-sm border border-slate-200 space-y-6">
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-1">
            Daily Connection Limit
          </label>
          <input
            type="number"
            data-testid="input-daily-connections"
            value={config.daily_connection_limit}
            onChange={(e) => setConfig({ ...config, daily_connection_limit: parseInt(e.target.value) || 0 })}
            className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
          />
          <p className="text-xs text-slate-500 mt-1">Recommended: 50-100 per day</p>
        </div>

        <div>
          <label className="block text-sm font-medium text-slate-700 mb-1">
            Daily Message Limit
          </label>
          <input
            type="number"
            data-testid="input-daily-messages"
            value={config.daily_message_limit}
            onChange={(e) => setConfig({ ...config, daily_message_limit: parseInt(e.target.value) || 0 })}
            className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Min Action Delay (ms)
            </label>
            <input
              type="number"
              value={config.min_action_delay_ms}
              onChange={(e) => setConfig({ ...config, min_action_delay_ms: parseInt(e.target.value) || 0 })}
              className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Max Action Delay (ms)
            </label>
            <input
              type="number"
              value={config.max_action_delay_ms}
              onChange={(e) => setConfig({ ...config, max_action_delay_ms: parseInt(e.target.value) || 0 })}
              className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Business Hours Start
            </label>
            <select
              value={config.business_hours_start}
              onChange={(e) => setConfig({ ...config, business_hours_start: parseInt(e.target.value) })}
              className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
            >
              {Array.from({ length: 24 }, (_, i) => (
                <option key={i} value={i}>{`${i.toString().padStart(2, "0")}:00`}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Business Hours End
            </label>
            <select
              value={config.business_hours_end}
              onChange={(e) => setConfig({ ...config, business_hours_end: parseInt(e.target.value) })}
              className="w-full px-4 py-2.5 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
            >
              {Array.from({ length: 24 }, (_, i) => (
                <option key={i} value={i}>{`${i.toString().padStart(2, "0")}:00`}</option>
              ))}
            </select>
          </div>
        </div>

        <div className="flex items-center justify-between pt-2">
          <span className="text-sm text-slate-700">Skip Weekends</span>
          <button
            onClick={() => setConfig({ ...config, skip_weekends: !config.skip_weekends })}
            className={`w-11 h-6 rounded-full transition-colors ${
              config.skip_weekends ? "bg-blue-600" : "bg-slate-200"
            }`}
          >
            <div
              className={`w-5 h-5 rounded-full bg-white shadow transform transition-transform ${
                config.skip_weekends ? "translate-x-5" : "translate-x-0.5"
              }`}
            />
          </button>
        </div>
      </div>
    </div>
  );
};

// Activity Logs Page
const ActivityLogsPage = () => {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState("");

  useEffect(() => {
    fetchLogs();
  }, [filter]);

  const fetchLogs = async () => {
    setLoading(true);
    try {
      const url = filter ? `/activity-logs?action_type=${filter}` : "/activity-logs";
      const response = await api.get(url);
      setLogs(response.data);
    } catch (error) {
      console.error("Failed to fetch logs", error);
    } finally {
      setLoading(false);
    }
  };

  const handleClear = async () => {
    if (!window.confirm("Clear all activity logs?")) return;
    try {
      await api.delete("/activity-logs");
      toast.success("Logs cleared");
      fetchLogs();
    } catch (error) {
      toast.error("Failed to clear logs");
    }
  };

  const statusIcon = (status) => {
    if (status === "success") return <CheckCircle className="w-4 h-4 text-green-500" />;
    if (status === "failure") return <X className="w-4 h-4 text-red-500" />;
    return <Info className="w-4 h-4 text-yellow-500" />;
  };

  return (
    <div data-testid="activity-logs-page">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Activity Logs</h2>
          <p className="text-slate-500">View automation activity history</p>
        </div>
        <div className="flex items-center gap-3">
          <select
            data-testid="filter-action-type"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
          >
            <option value="">All Actions</option>
            <option value="automation">Automation</option>
            <option value="credentials">Credentials</option>
            <option value="search">Search</option>
            <option value="connect">Connect</option>
            <option value="message">Message</option>
            <option value="config">Config</option>
            <option value="error">Error</option>
          </select>
          <button
            onClick={fetchLogs}
            className="p-2 hover:bg-slate-100 rounded-lg transition-colors"
          >
            <RefreshCw className="w-5 h-5 text-slate-600" />
          </button>
          <button
            onClick={handleClear}
            data-testid="btn-clear-logs"
            className="p-2 hover:bg-red-50 rounded-lg transition-colors"
          >
            <Trash2 className="w-5 h-5 text-red-500" />
          </button>
        </div>
      </div>

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <RefreshCw className="w-8 h-8 animate-spin text-blue-600" />
        </div>
      ) : logs.length === 0 ? (
        <div className="bg-white rounded-xl p-8 shadow-sm border border-slate-200 text-center">
          <Activity className="w-12 h-12 text-slate-300 mx-auto mb-3" />
          <p className="text-slate-500">No activity logs yet</p>
        </div>
      ) : (
        <div className="space-y-2">
          {logs.map((log) => (
            <div key={log.id} className="bg-white rounded-lg p-4 shadow-sm border border-slate-200 flex items-center gap-4">
              {statusIcon(log.status)}
              <div className="flex-1">
                <p className="font-medium text-slate-900">{log.description}</p>
                <p className="text-sm text-slate-500">
                  {log.action_type} • {new Date(log.timestamp).toLocaleString()}
                </p>
              </div>
              {log.details && (
                <button
                  onClick={() => toast.info(JSON.stringify(log.details, null, 2))}
                  className="text-sm text-blue-600 hover:underline"
                >
                  View Details
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

// Main App
function App() {
  const [activeTab, setActiveTab] = useState("dashboard");

  const renderPage = () => {
    switch (activeTab) {
      case "dashboard":
        return <DashboardPage />;
      case "credentials":
        return <CredentialsPage />;
      case "search":
        return <SearchCriteriaPage />;
      case "templates":
        return <TemplatesPage />;
      case "connections":
        return <ConnectionsPage />;
      case "messages":
        return <MessagesPage />;
      case "stealth":
        return <StealthConfigPage />;
      case "rate-limits":
        return <RateLimitsPage />;
      case "activity":
        return <ActivityLogsPage />;
      default:
        return <DashboardPage />;
    }
  };

  return (
    <div className="flex min-h-screen bg-slate-50" data-testid="app-container">
      <Toaster position="top-right" richColors />
      <Sidebar activeTab={activeTab} setActiveTab={setActiveTab} />
      <main className="flex-1 p-8">
        {renderPage()}
      </main>
    </div>
  );
}

export default App;
