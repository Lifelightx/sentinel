import { render } from 'preact';
import { html } from 'htm/preact';
import { useState, useEffect } from 'preact/hooks';

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
    const interval = setInterval(fetchServers, 5000);
    return () => clearInterval(interval);
  }, []);

  return html`
    <div class="panel">
      <table class="table">
        <thead>
          <tr>
            <th>Server</th>
            <th>Status</th>
            <th>CPU</th>
            <th>RAM</th>
            <th>Disk</th>
            <th>Alerts</th>
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
              <td>${new Date(server.LastSeen * 1000).toLocaleString()}</td>
            </tr>
          `)}
        </tbody>
      </table>
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

function ContainerList({ serverId }) {
  const [containers, setContainers] = useState([]);

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

  return html`
    <div class="panel">
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
          ${containers.length === 0 ? html`
            <tr><td colspan="7" class="empty">No containers found</td></tr>
          ` : containers.map(c => html`
            <tr>
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
              <td>${c.cpuPercentage.toFixed(1)}%</td>
              <td>${c.memoryMB.toFixed(1)} MB</td>
              <td>${c.memoryPercentage.toFixed(1)}%</td>
            </tr>
          `)}
        </tbody>
      </table>
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
  render(html`
    <${ServerSummary} serverId=${serverId} />
    <${ContainerList} serverId=${serverId} />
  `, serverAppEl);
}
