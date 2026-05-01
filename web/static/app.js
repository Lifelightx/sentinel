import { render } from 'preact';
import { html } from 'htm/preact';
import { useState, useEffect, useMemo } from 'preact/hooks';

const getThresholdColor = (val) => {
  if (val < 60) return 'green';
  if (val <= 80) return 'yellow';
  return 'red';
};

const sortPriority = (c) => {
  if (c.health === 'unhealthy') return 1;
  if (c.state === 'restarting') return 2;
  if (c.state === 'exited') return 3;
  if (c.state === 'created') return 4;
  if (c.health === 'healthy' && c.state === 'running') return 5;
  if (c.state === 'running') return 6;
  return 7;
};

function Toast({ message, onClose }) {
  useEffect(() => {
    if (!message) return;
    const timer = setTimeout(onClose, 3000);
    return () => clearTimeout(timer);
  }, [message, onClose]);
  
  if (!message) return null;
  return html`
    <div class="toast">
      ${message}
    </div>
  `;
}

function Modal({ isOpen, onClose, title, children }) {
  if (!isOpen) return null;
  return html`
    <div class="modal-overlay" onClick=${onClose}>
      <div class="modal-content" onClick=${e => e.stopPropagation()}>
        <div class="modal-header">
          <h2>${title}</h2>
          <button class="close-btn" onClick=${onClose}>✕</button>
        </div>
        <div class="modal-body">
          ${children}
        </div>
      </div>
    </div>
  `;
}

function Dashboard() {
  const [servers, setServers] = useState([]);

  useEffect(() => {
    const fetchServers = () => {
      fetch('/api/servers')
        .then(res => res.json())
        .then(data => setServers(data || []))
        .catch(console.error);
    };
    fetchServers();
    console.log(servers)
    const interval = setInterval(fetchServers, 5000);
    return () => clearInterval(interval);
  }, []);

  return html`
    <div class="panel">
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>Server</th>
              <th>Status</th>
              <th>CPU</th>
              <th>RAM</th>
              <th>Disk</th>
              <th>Alerts</th>
              <th>IP</th>
              <th>Last Seen</th>
            </tr>
          </thead>
          <tbody>
            ${servers.map(server => html`
              <tr class="row-link" onClick=${() => window.location='/server/'+server.ServerID}>
                <td>${server.ServerID}</td>
                <td>
                  ${server.Status === 'online' 
                    ? html`<span class="badge green">online</span>`
                    : html`<span class="badge red">offline</span>`}
                </td>
                <td>${server.CPU.toFixed(1)}%</td>
                <td>${server.RAM.toFixed(1)}%</td>
                <td>${server.Disk.toFixed(1)}%</td>
                <td>
                  ${server.AlertCount > 0 
                    ? html`<span class="badge red">${server.AlertCount}</span>`
                    : html`<span class="badge gray">0</span>`}
                </td>
                <td>${server.IPv4}</td>
                <td>${new Date(server.LastSeen * 1000).toLocaleString()}</td>
              </tr>
            `)}
          </tbody>
        </table>
      </div>
    </div>
  `;
}

function ServerSummary({ serverId }) {
  const [server, setServer] = useState(null);

  useEffect(() => {
    const fetchServer = () => {
      fetch('/api/servers/alerts?ServerID=' + serverId)
        .then(res => res.json())
        .then(data => setServer(data))
        .catch(console.error);
    };
    fetchServer();
    const interval = setInterval(fetchServer, 5000);
    return () => clearInterval(interval);
  }, [serverId]);

  if (!server) return html`<p>Loading summary...</p>`;

  return html`
    <div class="stats-bar">
      <div class="stat-box">
        <span class="label">Status</span>
        <span class="value ${server.Status === 'online' ? 'green-text' : 'red-text'}">
          ${server.Status}
        </span>
      </div>
      <div class="stat-box">
        <span class="label">CPU</span>
        <span class="value">${server.CPU.toFixed(1)}%</span>
      </div>
      <div class="stat-box">
        <span class="label">RAM</span>
        <span class="value">${server.RAM.toFixed(1)}%</span>
      </div>
      <div class="stat-box">
        <span class="label">Disk</span>
        <span class="value">${server.Disk.toFixed(1)}%</span>
      </div>
    </div>
  `;
}

function ContainerStats({ containers }) {
  const stats = {
    total: containers.length,
    healthy: containers.filter(c => c.health === 'healthy').length,
    warnings: containers.filter(c => c.health === 'unhealthy' || c.state === 'restarting').length,
    stopped: containers.filter(c => c.state === 'exited' || c.state === 'created').length,
  };

  return html`
    <div class="container-stats-bar">
      <div class="compact-stat">
        <span class="label">Containers</span>
        <span class="value">${stats.total}</span>
      </div>
      <div class="compact-stat">
        <span class="label">Healthy</span>
        <span class="value green-text">${stats.healthy}</span>
      </div>
      <div class="compact-stat">
        <span class="label">Warnings</span>
        <span class="value ${stats.warnings > 0 ? 'yellow-text' : ''}">${stats.warnings}</span>
      </div>
      <div class="compact-stat">
        <span class="label">Stopped</span>
        <span class="value muted">${stats.stopped}</span>
      </div>
    </div>
  `;
}

function ContainerList({ serverId, onSelectContainer }) {
  const [containers, setContainers] = useState([]);
  const [search, setSearch] = useState('');

  useEffect(() => {
    const fetchContainers = () => {
      fetch('/api/containers?ServerID=' + serverId)
        .then(res => res.json())
        .then(data => setContainers(data || []))
        .catch(console.error);
    };
    fetchContainers();
    const interval = setInterval(fetchContainers, 10000);
    return () => clearInterval(interval);
  }, [serverId]);

  const filteredAndSorted = useMemo(() => {
    let result = [...containers];
    if (search) {
      const q = search.toLowerCase();
      result = result.filter(c => c.name.toLowerCase().includes(q) || c.id.toLowerCase().includes(q));
    }
    result.sort((a, b) => sortPriority(a) - sortPriority(b));
    return result;
  }, [containers, search]);

  return html`
    <div>
      <${ContainerStats} containers=${containers} />
      
      <input 
        type="text" 
        class="search-input" 
        placeholder="Search containers..." 
        value=${search}
        onInput=${e => setSearch(e.target.value)}
      />

      <div class="panel">
        <div class="table-container">
          <table class="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>ID</th>
                <th>State</th>
                <th>Health</th>
                <th>CPU</th>
                <th>Memory</th>
                <th>Mem %</th>
              </tr>
            </thead>
            <tbody>
              ${filteredAndSorted.length === 0 ? html`
                <tr><td colspan="7" class="empty">No containers found</td></tr>
              ` : filteredAndSorted.map(c => html`
                <tr class="row-link" onClick=${() => onSelectContainer(c)}>
                  <td>${c.name}</td>
                  <td class="muted">${c.id.substring(0, 12)}</td>
                  <td>
                    ${c.state === 'running' ? html`<span class="badge green">running</span>`
                      : c.state === 'exited' ? html`<span class="badge red">exited</span>`
                      : c.state === 'created' ? html`<span class="badge yellow">created</span>`
                      : html`<span class="badge gray">${c.state}</span>`}
                  </td>
                  <td>
                    ${c.health === 'healthy' ? html`<span class="badge green">healthy</span>`
                      : c.health === 'unhealthy' ? html`<span class="badge red">unhealthy</span>`
                      : c.health === 'starting' ? html`<span class="badge yellow">starting</span>`
                      : html`<span class="badge gray">none</span>`}
                  </td>
                  <td>
                    <span class="${getThresholdColor(c.cpuPercentage)}-text">${c.cpuPercentage.toFixed(1)}%</span>
                  </td>
                  <td>
                    <span class="${getThresholdColor(c.memoryPercentage)}-text">${c.memoryMB.toFixed(1)} MB</span>
                  </td>
                  <td>
                    <span class="${getThresholdColor(c.memoryPercentage)}-text">${c.memoryPercentage.toFixed(1)}%</span>
                  </td>
                </tr>
              `)}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  `;
}

function ContainerDetails({ container, onBack }) {
  const [toastMsg, setToastMsg] = useState('');
  const [modalState, setModalState] = useState(null);
  const [modalContent, setModalContent] = useState("Loading...");
  const serverId = window.location.pathname.split("/")[2];

  const runAction = async (action) => {
    setModalState(action);
    setModalContent("Loading...");

    // trigger backend
    await fetch("/api/actions", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        action: action.toLowerCase(),
        container_id: container.id,
        hostName: serverId
      })
    });

    pollResult(action.toLowerCase());
  };

  const pollResult = (action) => {
    let attempts = 0;

    const interval = setInterval(async () => {
      attempts++;

      const res = await fetch(
        `/api/result?hostname=${serverId}&containerId=${container.id}&action=${action}`
      );

      if (res.status === 200) {
        const data = await res.json();
        clearInterval(interval);

        let content = data.Response.data;

        // pretty print JSON
        if (action === "inspect") {
          try {
            content = JSON.stringify(JSON.parse(content), null, 2);
          } catch {}
        }

        setModalContent(content);
      }

      if (attempts > 10) {
        clearInterval(interval);
        setModalContent("Timeout: no response from agent");
      }

    }, 1000);
  };

  return html`
    <div class="container-details">
      <div class="page-head" style="margin-bottom: 24px;">
        <div>
          <button class="btn" onClick=${onBack} style="margin-bottom: 12px;">← Back</button>
          <h1>${container.name}</h1>
          <div class="muted">ID: ${container.id}</div>
        </div>
      </div>

      <div class="panel">
        <div class="modal-header">
          <h3>Actions</h3>
          <div>
            <button class="action-btn" onClick=${() => runAction("logs")}>
              Logs
            </button>
            <button class="action-btn" onClick=${() => runAction("inspect")}>
              Inspect
            </button>
          </div>
        </div>
      </div>

      <${Modal}
        isOpen=${!!modalState}
        onClose=${() => setModalState(null)}
        title=${modalState === "logs" ? "Container Logs" : "Inspect"}
      >
        <div class=${modalState === "logs" ? "mock-logs" : "mock-json"}>
          ${modalContent}
        </div>
      </${Modal}>
    </div>
  `;
}

function ServerApp({ serverId }) {
  const [selectedContainer, setSelectedContainer] = useState(null);

  if (selectedContainer) {
    return html`<${ContainerDetails} container=${selectedContainer} onBack=${() => setSelectedContainer(null)} />`;
  }

  return html`
    <div>
      <${ServerSummary} serverId=${serverId} />
      <${ContainerList} serverId=${serverId} onSelectContainer=${setSelectedContainer} />
    </div>
  `;
}

const dashboardEl = document.getElementById('preact-dashboard');
if (dashboardEl) {
  render(html`<${Dashboard} />`, dashboardEl);
}

const serverAppEl = document.getElementById('preact-server-app');
if (serverAppEl) {
  const serverId = serverAppEl.dataset.serverId;
  render(html`<${ServerApp} serverId=${serverId} />`, serverAppEl);
}
