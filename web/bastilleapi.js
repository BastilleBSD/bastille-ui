function bastilleApp() {
    return {
        isAuthenticated: false,
        loading: false,
        theme: 'dark',
        currentView: 'jails',
        showTemplateInputModal: false,
        sidebarExpanded: window.innerWidth >= 768,
        showImportModal: false,
        showExportModal: false,
        bootstrapUpdate: true,
        showRefreshMenu: false,
        showCreateModal: false,
        autoRefreshTimer: null,
        viewMode: 'grid',
        templateFilter: '',
        insightsTab: 'ps',
        firewallData: {},
        servicesData: {},
        resourcesData: {},
        hostStats: { ncpu: 0, memory: 0, memoryHuman: '', storage: '' },
        showApiModal: false,
        showVerifyModal: false,
        showReleaseModal: false,
        showConsoleModal: false,
        showApplyModal: false,
        consoleUrl: '',
        newReleaseName: '',
        showFirewallModal: false,
        templateTargetJail: '',
        templateNameInput: '',
        consoleTitle: '',
        applyTargetTemplate: '',
        apiConnected: false,

        apiConfig: { url: '', key: '' },
        credentials: { username: '', password: '' },
        toast: { visible: false, message: '' },

        navItems: [
            { id: 'jails', label: 'Jails', icon: 'fa-solid fa-cubes' },
            { id: 'services', label: 'Services', icon: 'fa-solid fa-server' },
            { id: 'resources', label: 'Resources', icon: 'fa-solid fa-gauge-high' },
            { id: 'network', label: 'Network', icon: 'fa-solid fa-network-wired' },
            { id: 'storage', label: 'Storage', icon: 'fa-solid fa-hard-drive' },
            { id: 'releases', label: 'Releases', icon: 'fa-solid fa-tags' },
            { id: 'templates', label: 'Templates', icon: 'fa-solid fa-clone' },
            { id: 'firewall', label: 'Firewall', icon: 'fa-solid fa-fire' },
            { id: 'security', label: 'Security', icon: 'fa-solid fa-shield-halved' },
            { id: 'hardening', label: 'Hardening', icon: 'fa-solid fa-lock' },
            { id: 'insights', label: 'Insights', icon: 'fa-solid fa-magnifying-glass-chart' },
            { id: 'settings', label: 'Settings', icon: 'fa-solid fa-gear' },
        ],

        // Data Storage
        jails: [],
        releases: [],
        templates: [],
        backups: [],
        verifyOutput: '',
        selectedJails: [],
        networkMap: [],
        securityData: {},
        insightsData: {},
        hardeningData: {},
        
        // New Jail Form State
        newJail: {
            name: '',
            release: '',
            ip: '',
            interface: '', // Optional interface
            flags: {
                vnet: false,      // -V
                bridge: false,    // -B
                thick: false,     // -T
                clone: false,     // -C
                linux: false,     // -L
                empty: false,     // -E
                noBoot: false     // --no-boot
            }
        },
        // New Limit Form State
        newLimit: {
            target: '',
            resource: 'memoryuse', // memoryuse, pcpu, maxproc
            value: '',             // 1G, 50, etc.
            action: 'deny'         // deny, log, etc.
        },
        // Import Jail State
        newImport: {
            file: '',
            release: '',
            force: false,    // -f
            staticMac: false, // -M
            mac: ''           // Custom MAC address
        },
        // Export Jail State
        exportSettings: {
            target: '',
            compression: 'txz', // Default to standard archive
            live: false,        // --live
            auto: false         // --auto
        },
        // Network Management State
        networkForm: {
            target: '',
            action: 'add',      // add | remove
            interface: '',      // e.g. vtnet0, bridge0
            ip: '',             // e.g. 10.0.0.50/24 or DHCP
            vlan: '',           // VLAN ID (-v)
            flags: {
                auto: false,      // -a
                bridge: false,    // -B
                staticMac: false, // -M
                noIp: false,      // -n
                passthrough: false, // -P
                vnet: false       // -V
            }
        },
        // Concurrency Control
        requestQueue: [],
        activeRequests: 0,
        maxConcurrency: 2, // Limit simultaneous API calls to 3

        // New Firewall Rule State
        newRule: {
            target: '',
            protocol: 'tcp',
            hostPort: '',
            jailPort: ''
        },
        createTab: 'basic', // basic, network, options
        autoRefreshLabel: 'Off',

        // Mock data for static pages
        storage: [
            { name: 'zroot/bastille/jails/demo', used: '242M', avail: '6.23G', mountpoint: '/usr/local/bastille/jails/demo', compress: 'on', ratio: '1.27x' },
            { name: 'zroot/bastille/jails/demo/root', used: '241M', avail: '6.23G', mountpoint: '/usr/local/bastille/jails/demo/root', compress: 'on', ratio: '1.27x' }
        ],
        network: [
            { interface: 'vtnet0', type: 'Physical', ip: '192.168.1.150/24', options: 'uplink' },
            { interface: 'bastille0', type: 'Bridge', ip: '10.0.0.1/24', options: 'gateway' }
        ],
        config: {
            'bastille_prefix': '/usr/local/bastille',
            'bastille_zfs_enable': 'YES',
            'bastille_zfs_zpool': 'zroot',
            'bastille_network_loopback': 'bastille0'
        },

        init() {
            if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
                this.setTheme('dark');
            } else {
                this.setTheme('light');
            }

            const savedConfig = localStorage.getItem('bastilleApiConfig');
            if (savedConfig) {
                this.apiConfig = JSON.parse(savedConfig);
                if(this.apiConfig.url) this.apiConnected = true;
            }

            // Restore Session
            if (localStorage.getItem('bastilleAuth') === 'true') {
                this.isAuthenticated = true;
            }


            // Populate mock data if we have no API connection
            if (!this.apiConfig.url) {
                this.releases = ['13.2-RELEASE', '14.3-RELEASE', '15.0-CURRENT'];
                this.templates = [
                    { name: 'BastilleBSD/nginx', description: 'Nginx Web Server', source: 'GitLab' },
                    { name: 'BastilleBSD/postgresql', description: 'PostgreSQL 14', source: 'GitLab' }
                ];
            }
        },

        navigate(viewId) {
            this.currentView = viewId;
            // On mobile (width < 768px), auto-close sidebar after selection
            if (window.innerWidth < 768) {
                this.sidebarExpanded = false;
            }
        },


        // --- Queue Management ---
        async processQueue() {
            if (this.activeRequests >= this.maxConcurrency || this.requestQueue.length === 0) return;

            while (this.activeRequests < this.maxConcurrency && this.requestQueue.length > 0) {
                const task = this.requestQueue.shift();
                this.activeRequests++;
                
                // Execute task and ensure we decrement activeRequests even if it fails
                task().finally(() => {
                    this.activeRequests--;
                    this.processQueue();
                }).catch(err => {
                    console.error("Task failed:", err);
                    this.activeRequests--;
                    this.processQueue();
                });
            }
        },

        queueTask(taskFn) {
            this.requestQueue.push(taskFn);
            this.processQueue();
        },


        saveApiSettings() {
            if (!this.apiConfig.url) {
                this.showToast('API URL is required');
                return;
            }
            this.apiConfig.url = this.apiConfig.url.replace(/\/$/, "");

            localStorage.setItem('bastilleApiConfig', JSON.stringify(this.apiConfig));
            this.showApiModal = false;
            this.apiConnected = true;
            this.showToast('Settings saved. Refreshing...');
            this.refreshData();
        },

        login() {
            this.loading = true;
            
            // Hardcoded Credentials for Demo/MVP
            // In production, this should hit an auth endpoint
            setTimeout(() => {
                if(this.credentials.username === 'warden' && this.credentials.password === 'bastille') {
                    this.isAuthenticated = true;
                    localStorage.setItem('bastilleAuth', 'true');
                    this.showToast('Welcome back, ' + this.credentials.username);
                    this.refreshData();
                } else {
                    this.showToast('Invalid credentials. Try admin / bastille');
                }
                this.loading = false;
            }, 800);
        },

        logout() {
            this.isAuthenticated = false;
            localStorage.removeItem('bastilleAuth');
            this.credentials = { username: '', password: '' };
            this.currentView = 'jails';
        },

        toggleTheme() {
            const nextTheme = this.theme === 'light' ? 'dark' : 'light';
            this.setTheme(nextTheme);
        },

        setTheme(val) {
            this.theme = val;
            localStorage.theme = val;
            if (val === 'dark') {
                document.documentElement.classList.add('dark');
            } else {
                document.documentElement.classList.remove('dark');
            }
        },

        updateAutoRefresh(seconds, label) {
            this.autoRefreshLabel = label;
            this.showRefreshMenu = false;

            if (this.autoRefreshTimer) clearInterval(this.autoRefreshTimer);
            this.autoRefreshTimer = null;

            if (seconds > 0) {
                this.autoRefreshTimer = setInterval(() => this.refreshData(), seconds * 1000);
            }
        },

        get uniqueTags() {
            const tags = new Set();
            this.jails.forEach(j => {
                if (j.Tags && j.Tags !== '-') {
                    j.Tags.split(',').forEach(t => tags.add(t.trim()));
                }
            });
            return Array.from(tags).sort();
        },

        get filteredTemplates() {
            if (this.templateFilter === '') {
                return this.templates;
            }
            const lowerFilter = this.templateFilter.toLowerCase();
            return this.templates.filter(tpl => tpl.name.toLowerCase().includes(lowerFilter));
        },

        isTagSelected(tag) {
            const targetJails = this.jails.filter(j => j.Tags && j.Tags.split(',').includes(tag));
            return targetJails.length > 0 && targetJails.every(j => this.selectedJails.includes(j.Name));
        },

        selectJailsByTag(tag) {
            const targetJails = this.jails.filter(j => j.Tags && j.Tags.split(',').includes(tag));
            const allSelected = this.isTagSelected(tag);

            if (allSelected) {
                this.selectedJails = this.selectedJails.filter(name => !targetJails.find(j => j.Name === name));
            } else {
                targetJails.forEach(j => { if (!this.selectedJails.includes(j.Name)) this.selectedJails.push(j.Name); });
            }
        },

        parseStorageOutput(text) {
            const lines = text.split('\n');
            const datasets = [];

            lines.forEach(line => {
                line = line.trim();
                if (!line) return;

                // Skip headers (JAIL headers [xxx]: or ZFS headers NAME...)
                if (line.startsWith('[') || line.startsWith('NAME')) return;

                // Parse columns (whitespace separated)
                // NAME, USED, AVAIL, REFER, MOUNTPOINT, COMPRESS, RATIO
                const parts = line.split(/\s+/);
                if (parts.length >= 7) {
                    datasets.push({
                        name: parts[0],
                        used: parts[1],
                        avail: parts[2],
                        mountpoint: parts[4],
                        compress: parts[5],
                        ratio: parts[6]
                    });
                }
            });
            return datasets;
        },

        parseBytes(sizeStr) {
            if (!sizeStr) return 0;
            const units = { 'B': 1, 'K': 1024, 'M': 1024**2, 'G': 1024**3, 'T': 1024**4, 'P': 1024**5 };
            const match = sizeStr.toString().match(/^([\d.]+)([BKMGTPE]?)/i);
            if (!match) return 0;
            const unit = (match[2] || '').toUpperCase() || 'B';
            return parseFloat(match[1]) * (units[unit] || 1);
        },

        buildNetworkMap(text) {
            const lines = text.trim().split('\n');
            const interfaceMap = {};

            lines.forEach(line => {
                // Expected format: jailName: interface|ip.add.re.ss
                // However, the 'config' command might return just "interface|ip".
                // We handle the pipe split.
                if (line.includes('|')) {
                    const [iface, ip] = line.split('|');
                    const cleanIface = iface.trim();
                    const cleanIp = ip.trim();

                    if (!interfaceMap[cleanIface]) {
                        interfaceMap[cleanIface] = [];
                    }

                    // Find which jail owns this IP (cross-reference with this.jails)
                    const ownerJail = this.jails.find(j => j.IP === cleanIp);
                    const name = ownerJail ? ownerJail.Name : 'Unknown';

                    interfaceMap[cleanIface].push({
                        name: name,
                        ip: cleanIp,
                        status: ownerJail ? ownerJail.State : 'Down'
                    });
                }
            });

            // Convert map to array for Alpine x-for
            this.networkMap = Object.keys(interfaceMap).map(iface => ({
                name: iface,
                // Heuristic: bastille0 is always NAT/Loopback. Others (em0, vtnet0) are Bridged.
                type: iface === 'bastille0' ? 'NAT' : 'BRIDGE',
                jails: interfaceMap[iface]
            }));
        },

        openExportModal(jailName) {
            this.exportSettings.target = jailName;
            this.exportSettings.compression = 'txz';
            this.exportSettings.live = false;
            this.exportSettings.auto = false;
            this.showExportModal = true;
        },

        async performExport() {
            const { target, compression, live, auto } = this.exportSettings;
            this.showExportModal = false;
            this.showToast(`Exporting ${target}...`);

            try {
                let args = `--${compression}`;
                if (live) args += ' --live';
                if (auto) args += ' --auto';

                const endpoint = `${this.apiConfig.url}/api/v1/bastille/export?options=${encodeURIComponent(args.trim())}&target=${target}`;
                await this.fetchWithAuth(endpoint, 'POST');
                this.showToast(`Export started for ${target}. Check /usr/local/bastille/backups.`);
            } catch (e) {
                this.showToast(`Export Failed: ${e.message}`);
            }
        },

        async snapshotJail(jailName) {
            // Create ISO timestamp tag (e.g. 2023-10-27T10-30-00Z)
            const tag = new Date().toISOString().replace(/:/g, '-');
            this.showToast(`Snapshotting ${jailName}@${tag}...`);

            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/zfs?target=${jailName}&action=snapshot&tag=${tag}`;
                await this.fetchWithAuth(endpoint, 'POST');
                this.showToast(`Snapshot created for ${jailName}`);
            } catch (e) {
                this.showToast(`Snapshot Failed: ${e.message}`);
            }
        },

        async rollbackJail(jailName) {
            if (!confirm(`Are you sure you want to rollback ${jailName}? This will revert to the last snapshot and cannot be undone.`)) return;
            this.showToast(`Rolling back ${jailName}...`);

            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/zfs?target=${jailName}&action=rollback`;
                await this.fetchWithAuth(endpoint, 'POST');
                this.showToast(`${jailName} rolled back successfully.`);
            } catch (e) {
                this.showToast(`Rollback Failed: ${e.message}`);
            }
        },

        async verifyTemplate(templateName) {
            this.showToast(`Inspecting ${templateName}...`);
            if (!this.apiConfig.url) return;

            try {
                // POST /api/v1/bastille/verify?target=path/to/template
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/verify?target=${templateName}`;
                const headers = {};
                if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                headers['Authorization-ID'] = `bastille`;

                const response = await fetch(endpoint, { method: 'POST', headers: headers });

                if (!response.ok) throw new Error(`API Error: ${response.status}`);

                const text = await response.text();
                this.verifyOutput = text.split('\n').filter(line => !line.includes('Detected Bastillefile hook.')).join('\n');
                this.showVerifyModal = true;

            } catch (error) {
                this.showToast(`Verification Failed: ${error.message}`);
            }
        },

        closeConsole() {
            this.showConsoleModal = false;
            this.consoleUrl = ''; // Clear src to kill connection
        },

        async bootstrapRelease() {
            if (!this.newReleaseName) return;
            this.loading = true; // Lock UI
            this.showToast(`Bootstrapping ${this.newReleaseName}... this may take a while.`);

            try {
                let endpoint = `${this.apiConfig.url}/api/v1/bastille/bootstrap?target=${this.newReleaseName}`;
                if (this.bootstrapUpdate) {
                    endpoint = `${this.apiConfig.url}/api/v1/bastille/bootstrap?option=-u&target=${this.newReleaseName}`;
                }

                const headers = {};
                if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                headers['Authorization-ID'] = `bastille`;

                const response = await fetch(endpoint, { method: 'POST', headers: headers });

                if (!response.ok) throw new Error(`Bootstrap Failed: ${response.status}`);

                this.showToast(`${this.newReleaseName} fetched successfully.`);
                this.newReleaseName = '';
                this.showReleaseModal = false;
                await this.refreshData();

            } catch (error) {
                this.showToast(`Error: ${error.message}`);
            } finally {
                this.loading = false;
            }
        },

        async pullTemplates() {
            this.showToast('Pulling latest templates...');

            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/bootstrap?target=https://github.com/bastillebsd/templates`;
                const headers = {};
                if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                headers['Authorization-ID'] = `bastille`;

                const response = await fetch(endpoint, { method: 'POST', headers: headers });
                if (!response.ok) throw new Error(`Pull Failed: ${response.status}`);

                this.showToast('Templates updated successfully.');
                await this.refreshData();
            } catch (error) {
                this.showToast(`Error: ${error.message}`);
            }
        },

        async updateRelease(releaseName) {
            if (!this.apiConfig.url) return;
            this.showToast(`Updating ${releaseName}... this may take a while.`);
            
            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/update?target=${releaseName}`;
                await this.fetchWithAuth(endpoint, 'POST');
                
                this.showToast(`${releaseName} updated successfully.`);
            } catch (error) {
                this.showToast(`Update Failed: ${error.message}`);
            }
        },


        async destroyRelease(releaseName) {
            if (!this.apiConfig.url) return;
            if (!confirm(`Are you sure you want to destroy ${releaseName}? This cannot be undone.`)) return;

            this.showToast(`Destroying ${releaseName}...`);

            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/destroy?target=${releaseName}`;
                const headers = {};
                if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                headers['Authorization-ID'] = `bastille`;

                const response = await fetch(endpoint, { method: 'POST', headers: headers });
                if (!response.ok) throw new Error(`Destroy Failed: ${response.status}`);

                this.showToast(`${releaseName} destroyed.`);
                await this.refreshData();
            } catch (error) {
                this.showToast(`Error: ${error.message}`);
            }
        },

        async updateJailConfig(jailName, property, value) {
            this.showToast(`Updating ${property} to ${value}...`);

            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/config?target=${jailName}&action=set&property=${property}&value=${value}`;
                await this.fetchWithAuth(endpoint, 'POST');

                this.showToast(`Successfully set ${property} to ${value}`);
                // Slight delay to allow backend to write config before refreshing UI
                setTimeout(() => this.refreshData(), 500);
            } catch (e) {
                console.error(e);
                this.showToast(`Update Failed: ${e.message}`);
            }
        },


        async apiAction(jailName, action) {
            // Safety check for destructive actions
            if (action === 'destroy') {
                if (!confirm(`WARNING: Are you sure you want to DESTROY ${jailName}?\n\nThis action will permanently delete the jail and all its data.\nThis cannot be undone.`)) return;
            }

            this.showToast(`Command sent: ${action} on ${jailName}`);

            // Handle Console specially
            if (action === 'console') {
                await this.launchConsole(jailName);
                return;
            }


            if (this.apiConfig.url) {
                try {
                    const endpoint = `${this.apiConfig.url}/api/v1/bastille/${action}?target=${jailName}`;
                    const headers = {};
                    if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                    headers['Authorization-ID'] = `bastille`;

                    await fetch(endpoint, { method: 'POST', headers: headers });
                    setTimeout(() => this.refreshData(), 2000); // Refresh to confirm state change
                } catch (e) { console.error(e); }
            }

            const jail = this.jails.find(j => j.Name === jailName);
            if(jail) {
                if(action === 'stop') {
                    jail.State = 'Down';
                    jail.JID = '0';
                }
                if(action === 'destroy') {
                    // Remove the jail from the local list immediately
                    this.jails = this.jails.filter(j => j.Name !== jailName);
                }
                if(action === 'start' || action === 'restart') {
                    jail.State = 'Up';
                }
            }
        },

        async launchConsole(jailName) {
            if (!this.apiConfig.url) return;

            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/console?target=${jailName}`;
                const headers = {};
                if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                headers['Authorization-ID'] = `bastille`;

                // 1. Call API to spawn the process
                const response = await fetch(endpoint, { headers: headers });

                if (!response.ok) throw new Error(`API Error: ${response.status}`);

                // Check for dynamic port header (X-TTYD-Port)
                const port = response.headers.get("X-TTYD-Port");

                if (port) {
                    // If the API gave us a specific port, use it.
                    // Note: This constructs a direct URL (e.g. https://host:8183).
                    // This requires the port to be open on the firewall.
                    const urlObj = new URL(this.apiConfig.url);
                    this.consoleUrl = `${urlObj.protocol}//${urlObj.hostname}:${port}/`;
                } else {
                    // Fallback to the static Nginx proxy path (defaulting to 8182 internally)
                    this.consoleUrl = `${this.apiConfig.url}/console/`;
                }

                this.consoleTitle = `root@${jailName} ~ # `;
                this.showConsoleModal = true;

            } catch (error) {
                this.showToast(`Console Failed: ${error.message}`);
            }
        },

        openApplyModal(templateName) {
            this.applyTargetTemplate = templateName;
            this.selectedJails = []; // Reset selection
            this.showApplyModal = true;
        },

        toggleJailSelection(jailName) {
            if (this.selectedJails.includes(jailName)) {
                this.selectedJails = this.selectedJails.filter(name => name !== jailName);
            } else {
                this.selectedJails.push(jailName);
            }
        },

        toggleSelectAll() {
            if (this.selectedJails.length === this.jails.length) {
                this.selectedJails = [];
            } else {
                this.selectedJails = this.jails.map(j => j.Name);
            }
        },

        openTemplateInput(jailName) {
            this.templateTargetJail = jailName;
            this.templateNameInput = '';
            this.showTemplateInputModal = true;
        },

        openCreateModal() {
            // Reset Form
            this.newJail = {
                name: '',
                release: this.releases.length > 0 ? this.releases[0] : '', // Default to first release
                ip: '',
                interface: '',
                flags: {
                    vnet: false, bridge: false, thick: false, 
                    clone: false, linux: false, empty: false, noBoot: false
                }
            };
            this.createTab = 'basic';
            this.showCreateModal = true;
        },

        async createJail() {
            // Validation
            if (!this.newJail.name || !this.newJail.release || !this.newJail.ip) {
                this.showToast('Name, Release, and IP are required.');
                return;
            }

            this.loading = true;
            this.showCreateModal = false;
            this.showToast(`Creating jail ${this.newJail.name}...`);

            // Build Options String
            let options = [];
            if (this.newJail.flags.vnet) options.push('-V');
            if (this.newJail.flags.bridge) options.push('-B');
            if (this.newJail.flags.thick) options.push('-T');
            if (this.newJail.flags.clone) options.push('-C');
            if (this.newJail.flags.linux) options.push('-L');
            if (this.newJail.flags.empty) options.push('-E');
            if (this.newJail.flags.noBoot) options.push('--no-boot');

            try {
                // Construct API URL: /api/v1/bastille/create?name=X&release=Y&ip=Z&iface=W&options=...
                const params = new URLSearchParams();

                // Strip patch suffix (e.g. -p12) for the API call
                const cleanRelease = this.newJail.release.replace(/-p\d+$/, '');

                params.append('name', this.newJail.name);
                params.append('release', cleanRelease);
                params.append('ip', this.newJail.ip);
                if (this.newJail.interface) params.append('iface', this.newJail.interface);
                if (options.length > 0) params.append('options', options.join(' '));

                const endpoint = `${this.apiConfig.url}/api/v1/bastille/create?${params.toString()}`;
                await this.fetchWithAuth(endpoint, 'POST');
                
                this.showToast(`Jail ${this.newJail.name} created successfully.`);
                await this.refreshData();
            } catch (e) {
                this.showToast(`Creation Failed: ${e.message}`);
            } finally {
                this.loading = false;
            }
        },


        async submitTemplateInput() {
            // Validate syntax: author/template (alphanumeric)
            const regex = /^[a-z0-9]+\/[a-z0-9\-_]+$/i;
            if (!regex.test(this.templateNameInput)) {
                this.showToast("Invalid format. Use 'author/template'");
                return;
            }

            this.showTemplateInputModal = false;
            this.showToast(`Applying ${this.templateNameInput} to ${this.templateTargetJail}...`);

            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/template?target=${this.templateTargetJail}&arg=${this.templateNameInput}`;
                await this.fetchWithAuth(endpoint, 'POST');
                this.showToast(`Template applied to ${this.templateTargetJail}`);
            } catch (e) {
                this.showToast(`Failed: ${e.message}`);
            }
        },

        async executeTemplate() {
            this.showApplyModal = false;
            this.showToast(`Applying ${this.applyTargetTemplate} to ${this.selectedJails.length} jails...`);

            if (this.apiConfig.url) {
                const headers = {};
                if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                headers['Authorization-ID'] = `bastille`;

                for (const jail of this.selectedJails) {
                    try {
                        // POST /api/v1/bastille/template?target=JAIL&arg=TEMPLATE
                        const endpoint = `${this.apiConfig.url}/api/v1/bastille/template?target=${jail}&arg=${this.applyTargetTemplate}`;
                        await fetch(endpoint, { method: 'POST', headers: headers });
                        this.showToast(`Applied to ${jail}`);
                    } catch (e) {
                        console.error(e);
                        this.showToast(`Failed to apply to ${jail}`);
                    }
                }
            }
        },

        async scanSecurity(targetJail = null) {
            const targets = targetJail ? [targetJail] : this.jails.map(j => j.Name);

            targets.forEach(name => {
                if (!this.securityData[name]) this.securityData[name] = {};
                this.securityData[name].loading = true;
                this.securityData[name].output = '';
            });

            const headers = {};
            if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
            headers['Authorization-ID'] = `bastille`;

            for (const name of targets) {
                try {
                    const endpoint = `${this.apiConfig.url}/api/v1/bastille/pkg?target=${name}&args=audit -F`;
                    
                    // USE RAW FETCH: pkg audit returns 500 (exit code 1) if vulns found.
                    // We must read the text body regardless of status.
                    const response = await fetch(endpoint, { method: 'POST', headers: headers });
                    const text = await response.text();

                    const analysis = this.parseSecurityOutput(text);

                    this.securityData[name] = {
                        loading: false,
                        output: text,
                        status: analysis.status,
                        count: analysis.count
                    };
                } catch (e) {
                    this.securityData[name] = { loading: false, status: 'error', output: e.message, count: 0 };
                }
            }
        },

        async upgradePkgs(targetJail = null) {
            const targets = targetJail ? [targetJail] : this.jails.map(j => j.Name);

            targets.forEach(name => {
                if (!this.securityData[name]) this.securityData[name] = {};
                this.securityData[name].loading = true;
                this.securityData[name].output = 'Starting package upgrade...';
                this.securityData[name].status = 'working';
            });

            const headers = {};
            if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
            headers['Authorization-ID'] = `bastille`;

            for (const name of targets) {
                try {
                    const endpoint = `${this.apiConfig.url}/api/v1/bastille/pkg?target=${name}&args=upgrade -y`;
                    
                    // USE RAW FETCH
                    const response = await fetch(endpoint, { method: 'POST', headers: headers });
                    const text = await response.text();

                    this.securityData[name] = {
                        loading: false,
                        output: text,
                        status: 'upgraded',
                        count: 0
                    };
                } catch (e) {
                    this.securityData[name] = { loading: false, status: 'error', output: e.message, count: 0 };
                }
            }
        },

        async scanHardening(targetJail = null) {
            const targets = targetJail ? [targetJail] : this.jails.map(j => j.Name);
            const headers = {};
            if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
            headers['Authorization-ID'] = `bastille`;

            targets.forEach(name => {
                this.hardeningData[name] = { loading: true, score: 0, details: {} };
            });

            targets.forEach(name => {
                this.queueTask(async () => {
                    try {
                    // 1. PKG Audit (Catch 500 error for vulns)
                    const pPkg = fetch(`${this.apiConfig.url}/api/v1/bastille/pkg?target=${name}&args=audit -F`, { method: 'POST', headers }).then(r => r.text());
                    
                    // 2. Open Ports (RDR list)
                    const pRdr = fetch(`${this.apiConfig.url}/api/v1/bastille/rdr?target=${name}&action=list`, { method: 'POST', headers }).then(r => r.text());
                    
                    // 3. Securelevel
                    const pSecure = fetch(`${this.apiConfig.url}/api/v1/bastille/config?target=${name}&action=get&property=securelevel`, { method: 'POST', headers }).then(r => r.text());
                    
                    // 4. Raw Sockets
                    const pRaw = fetch(`${this.apiConfig.url}/api/v1/bastille/config?target=${name}&action=get&property=allow.raw_sockets`, { method: 'POST', headers }).then(r => r.text());

                    // 5. Isolation (allow.mount, allow.sysvipc)
                    const pMount = fetch(`${this.apiConfig.url}/api/v1/bastille/config?target=${name}&action=get&property=allow.mount`, { method: 'POST', headers }).then(r => r.text());
                    const pIpc = fetch(`${this.apiConfig.url}/api/v1/bastille/config?target=${name}&action=get&property=allow.sysvipc`, { method: 'POST', headers }).then(r => r.text());

                    // 6. Services Audit (SSH & Sendmail)
                    const pSshRc = fetch(`${this.apiConfig.url}/api/v1/bastille/sysrc?target=${name}&args=sshd_enable`, { method: 'POST', headers }).then(r => r.text());
                    const pSshCfg = fetch(`${this.apiConfig.url}/api/v1/bastille/cmd?target=${name}&args=grep -i "^PermitRootLogin" /etc/ssh/sshd_config`, { method: 'POST', headers }).then(r => r.text());
                    const pSendmail = fetch(`${this.apiConfig.url}/api/v1/bastille/sysrc?target=${name}&args=sendmail_enable`, { method: 'POST', headers }).then(r => r.text());

                    const [resPkg, resRdr, resSecure, resRaw, resMount, resIpc, resSshRc, resSshCfg, resSendmail] = await Promise.all([pPkg, pRdr, pSecure, pRaw, pMount, pIpc, pSshRc, pSshCfg, pSendmail]);

                    // --- SCORING ALGORITHM ---
                    let score = 0;
                    let details = {};

                    // Metric 1: Vulnerabilities (Max 30)
                    const vulnCount = (resPkg.match(/is vulnerable:/g) || []).length;
                    let scorePkg = 30;
                    if (vulnCount > 0) scorePkg = 20;
                    if (vulnCount > 2) scorePkg = 10;
                    if (vulnCount > 5) scorePkg = 0;
                    score += scorePkg;
                    details.vulns = vulnCount;

                    // Metric 2: Securelevel (Max 30)
                    // Output format is usually "jailname: value" or just "value". We look for the number.
                    const secureLevel = parseInt(resSecure.replace(/[^0-9-]/g, '')) || 0; // Default to 0 if parse fails
                    let scoreSl = 5; // Default level 0
                    if (secureLevel >= 2) scoreSl = 20;
                    else if (secureLevel === 1) scoreSl = 10;
                    else if (secureLevel < 0) scoreSl = -20;
                    score += scoreSl;
                    details.securelevel = secureLevel;

                    // Metric 3: Raw Sockets (Max 10)
                    const rawSockets = parseInt(resRaw.replace(/[^0-9]/g, '')) || 0;
                    let scoreRaw = (rawSockets === 0) ? 10 : 0;
                    score += scoreRaw;
                    details.rawSockets = rawSockets;

                    // Metric 4: Isolation (Max 15)
                    const allowMount = parseInt(resMount.replace(/[^0-9]/g, '')) || 0;
                    const allowIpc = parseInt(resIpc.replace(/[^0-9]/g, '')) || 0;
                    let scoreIso = (allowMount === 0 && allowIpc === 0) ? 15 : 0;
                    score += scoreIso;
                    details.isolation = { mount: allowMount, ipc: allowIpc };

                    // Metric 5: SSH Audit (Max 10)
                    const sshEnabled = resSshRc.includes('YES');
                    const rootLogin = resSshCfg.toLowerCase().includes('permitrootlogin yes');
                    let scoreSsh = 5;
                    if (sshEnabled) {
                        if (rootLogin) scoreSsh = 0; // High risk
                        else scoreSsh = 10; // Secure SSH
                    }
                    score += scoreSsh;
                    details.ssh = { enabled: sshEnabled, rootLogin: rootLogin };

                    // Metric 6: Sendmail (Max 10)
                    // Default is often "NO" or "NONE" for secure jails
                    const sendmailEnabled = !resSendmail.includes('NO') && !resSendmail.includes('NONE');
                    let scoreSendmail = sendmailEnabled ? 0 : 5;
                    score += scoreSendmail;
                    details.sendmail = sendmailEnabled;


                    const rdrLines = resRdr.trim().split('\n').filter(l => l.trim().length > 0);
                    let scorePorts = 10;
                    let portDisplay = [];
                    let isVNET = false;

                    if (resRdr.includes('VNET jails do not support rdr')) {
                        isVNET = true;
                    } else {
                        // Parse PF output: proto tcp ... port = 8080 -> ... port 8080
                        rdrLines.forEach(line => {
                            const match = line.match(/proto\s+(tcp|udp).*?port\s*=\s*(\d+).*?->.*?port\s+(\d+)/);
                            if (match) {
                                portDisplay.push(`${match[1]}/${match[2]}:${match[3]}`);
                            }
                        });

                        if (portDisplay.length > 0) scorePorts = 5;
                    }

                    score += scorePorts;
                    details.openPorts = portDisplay;
                    details.isVNET = isVNET;

                    // Compound Penalty: Exposed Vulnerabilities
                    if (vulnCount > 0 && portDisplay.length > 0 && !isVNET) {
                        score -= 15;
                    }

                    // Clamp score 0-100
                    score = Math.max(0, Math.min(100, score));

                    this.hardeningData[name] = { loading: false, score, details };

                } catch (e) {
                        console.error(e);
                        this.hardeningData[name] = { loading: false, score: 0, error: true, details: {} };
                    }
                });
            });
        },

        openFirewallModal(jailName) {
            this.newRule.target = jailName;
            this.newRule.hostPort = '';
            this.newRule.jailPort = '';
            this.newRule.protocol = 'tcp';
            this.showFirewallModal = true;
        },

        async submitFirewallRule() {
            if (!this.newRule.hostPort || !this.newRule.jailPort) {
                this.showToast('Host and Jail ports are required.');
                return;
            }

            this.showFirewallModal = false;
            this.showToast(`Adding rule to ${this.newRule.target}...`);

            try {
                // /api/v1/bastille/rdr?target=JAIL&protocol=TCP&host_port=X&jail_port=Y
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/rdr?target=${this.newRule.target}&protocol=${this.newRule.protocol}&host_port=${this.newRule.hostPort}&jail_port=${this.newRule.jailPort}`;
                await this.fetchWithAuth(endpoint, 'POST');
                
                this.showToast('Rule added successfully.');
                // Slight delay to allow pf to reload before rescanning
                setTimeout(() => this.scanFirewall(this.newRule.target), 1000);
            } catch (e) {
                this.showToast(`Error adding rule: ${e.message}`);
            }
        },

        async firewallAction(jailName, action) {
            if (!confirm(`Are you sure you want to ${action.toUpperCase()} rules for ${jailName}?`)) return;
            
            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/rdr?target=${jailName}&action=${action}`;
                await this.fetchWithAuth(endpoint, 'POST');
                this.showToast(`Rules ${action}ed for ${jailName}.`);
                setTimeout(() => this.scanFirewall(jailName), 1000);
            } catch (e) {
                this.showToast(`Action failed: ${e.message}`);
            }
        },

        async scanFirewall(targetJail = null) {
            const targets = targetJail ? [targetJail] : this.jails.map(j => j.Name);
            const headers = {};
            if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
            headers['Authorization-ID'] = `bastille`;

            targets.forEach(name => {
                if (!this.firewallData[name]) this.firewallData[name] = { rules: [], loading: true, isVNET: false };
                else this.firewallData[name].loading = true;
            });

            targets.forEach(name => {
                this.queueTask(async () => {
                    try {
                    const endpoint = `${this.apiConfig.url}/api/v1/bastille/rdr?target=${name}&action=list`;
                    const response = await fetch(endpoint, { method: 'POST', headers });
                    const text = await response.text();

                    // Check for VNET exclusion
                    if (text.includes('VNET jails do not support rdr')) {
                        this.firewallData[name] = { rules: [], loading: false, isVNET: true };
                        return;
                    }

                    // Parse Rules
                    // Format: rdr pass on em0 inet proto tcp from any to any port = 8080 -> 10.0.0.80 port 8080
                    const rules = [];
                    const lines = text.trim().split('\n');
                    lines.forEach(line => {
                        if (!line.startsWith('rdr')) return;
                        
                        // Regex to extract key fields
                        const match = line.match(/on\s+(\S+)\s+(?:inet|inet6)\s+proto\s+(tcp|udp).*?port\s*=\s*(\d+).*?->\s*(\S+)\s+port\s+(\d+)/);
                        if (match) {
                            rules.push({
                                interface: match[1],
                                protocol: match[2],
                                hostPort: match[3],
                                jailIp: match[4],
                                jailPort: match[5],
                                raw: line
                            });
                        }
                    });

                    this.firewallData[name] = { rules: rules, loading: false, isVNET: false };

                    } catch (e) {
                        this.firewallData[name] = { rules: [], loading: false, isVNET: false, error: e.message };
                    }
                });
            });
        },

        async scanServices() {
            this.showToast('Scanning services (sysrc)...');
            // Initialize state for all known jails
            this.jails.forEach(j => {
                if (!this.servicesData[j.Name]) {
                    this.servicesData[j.Name] = { services: [], loading: true };
                } else {
                    this.servicesData[j.Name].loading = true;
                }
            });

            try {
                // Fetch all effective rc.conf variables
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/sysrc?target=ALL&args=-ae`;
                const response = await this.fetchWithAuth(endpoint, 'POST');
                const text = await response.text();

                // Parse output: jailname: service_suffix: value
                // e.g. "myjail: sshd_enable: YES"
                const serviceMap = {}; // { jail: { serviceName: { enable: 'NO', flags: '' } } }

                const lines = text.trim().split('\n');
                let currentJail = null;
                
                // Regex for Jail Header: [jailname]:
                const jailRegex = /^\[([a-zA-Z0-9\-\._]+)\]:$/;
                // Regex for Config Line: service_suffix="value"
                const configRegex = /^([a-zA-Z0-9_]+)_(enable|flags)="?(.*?)"?$/;

                const ignoredServices = ['sendmail_msp_queue', 'sendmail_outbound', 'sendmail_submit'];

                lines.forEach(line => {
                    line = line.trim();
                    if (!line) return;

                    // 1. Detect Jail Header
                    const jailMatch = line.match(jailRegex);
                    if (jailMatch) {
                        currentJail = jailMatch[1];
                        if (!serviceMap[currentJail]) serviceMap[currentJail] = {};
                        return;
                    }

                    // 2. Parse Config Lines
                    if (currentJail) {
                        const configMatch = line.match(configRegex);
                        if (configMatch) {
                            const serviceName = configMatch[1];
                            const suffix = configMatch[2]; // enable or flags
                            const value = configMatch[3];

                            if (ignoredServices.includes(serviceName)) return;

                            if (!serviceMap[currentJail][serviceName]) {
                                serviceMap[currentJail][serviceName] = { 
                                    name: serviceName, 
                                    enable: 'NO', // Default assumption
                                    flags: '', 
                                    statusOutput: 'Pending...', // Initialize tooltip
                                    status: 'Unknown', 
                                    pid: '', 
                                    running: false, 
                                    loading: false 
                                };
                            }

                            if (suffix === 'enable') serviceMap[currentJail][serviceName].enable = value;
                            if (suffix === 'flags') serviceMap[currentJail][serviceName].flags = value;
                        }
                    }
                });

                // Convert map to array and trigger status checks
                for (const jailName of Object.keys(serviceMap)) {
                    const list = Object.values(serviceMap[jailName]);
                    list.forEach(svc => {
                        svc.isEnabled = (svc.enable.toUpperCase() === 'YES') || (svc.flags.length > 0 && svc.enable.toUpperCase() !== 'NO');
                    });
                    
                    this.servicesData[jailName] = { services: list, loading: false };

                    this.servicesData[jailName].services.forEach(svc => {
                        this.queueTask(() => this.checkServiceStatus(jailName, svc));
                    });
                }

                // Clear loading for jails with no services found
                this.jails.forEach(j => {
                    if (this.servicesData[j.Name] && this.servicesData[j.Name].loading) {
                        this.servicesData[j.Name].loading = false;
                    }
                });

            } catch (e) {
                console.error(e);
                this.showToast(`Service Scan Failed: ${e.message}`);
            }
        },

        async checkServiceStatus(jailName, serviceObj) {
            // Updates the specific service object in place
            serviceObj.loading = true;
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 5000); // 5s timeout
            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/service?target=${jailName}&service=${serviceObj.name}&args=status`;
                // Use raw fetch to handle potential non-200 exit codes from 'service status'
                const headers = {};
                if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                headers['Authorization-ID'] = `bastille`;
                
                const response = await fetch(endpoint, { method: 'POST', headers: headers, signal: controller.signal });
                clearTimeout(timeoutId);

                const text = await response.text();

                serviceObj.statusOutput = text;
                // Handle 500 as "Stopped" (Exit Code 1), 200 as "Running" (Exit Code 0)
                if (response.ok) {
                    serviceObj.running = true;
                    // Extract PID if running (e.g. "cron is running as pid 11488")
                    const pidMatch = text.match(/pid\s+(\d+)/i);
                    serviceObj.pid = pidMatch ? pidMatch[1] : '';
                } else if (response.status === 500) {
                    serviceObj.running = false;
                    serviceObj.pid = '';
                    // Ensure we have text for the tooltip even if body is empty
                    if (!serviceObj.statusOutput || serviceObj.statusOutput.trim() === '') {
                        serviceObj.statusOutput = 'Service is stopped.';
                    }
                } else {
                    // Actual API error (404, 403, etc)
                    serviceObj.running = false;
                    serviceObj.statusOutput = `API Error: ${response.status}`;
                }

            } catch (e) {
                if (e.name === 'AbortError') serviceObj.statusOutput = 'Request timed out';
                else serviceObj.statusOutput = `Error: ${e.message}`;
                serviceObj.running = false;
            } finally {
                serviceObj.loading = false;
            }
        },

        async controlService(jailName, serviceName, action) {
            this.showToast(`${action.toUpperCase()} ${serviceName} on ${jailName}...`);
            const targetSvc = this.servicesData[jailName].services.find(s => s.name === serviceName);
            if(targetSvc) targetSvc.loading = true;

            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/service?target=${jailName}&service=${serviceName}&args=${action}`;
                await this.fetchWithAuth(endpoint, 'POST');
                
                // Re-check status after action
                setTimeout(() => this.checkServiceStatus(jailName, targetSvc), 1000);
                // If we toggled enable/disable, we should ideally re-scan sysrc, but for now we rely on the action succeeding.
                if(action === 'enable' || action === 'disable') {
                    if(targetSvc) targetSvc.isEnabled = (action === 'enable');
                }
            } catch (e) {
                this.showToast(`Action Failed: ${e.message}`);
                if(targetSvc) targetSvc.loading = false;
            }
        },

        async fetchHostStats() {
            // We only need to check one jail to get host hardware stats usually
            if(this.jails.length === 0) return;
            const target = this.jails[0].Name;

            try {
                // Fetch CPU Count
                const urlCpu = `${this.apiConfig.url}/api/v1/bastille/cmd?target=${target}&command=sysctl -n hw.ncpu`;
                const resCpu = await this.fetchWithAuth(urlCpu, 'POST');
                const textCpu = await resCpu.text();
                // Strip jail header (e.g. "[bomb]: 6" -> "6")
                const cpuMatch = textCpu.match(/:?\s*(\d+)\s*$/);
                this.hostStats.ncpu = cpuMatch ? cpuMatch[1] : '-';

                // Fetch Physical Memory
                const urlMem = `${this.apiConfig.url}/api/v1/bastille/cmd?target=${target}&command=sysctl -n hw.physmem`;
                const resMem = await this.fetchWithAuth(urlMem, 'POST');
                const textMem = await resMem.text();
                const memMatch = textMem.match(/:?\s*(\d+)\s*$/);
                const bytes = memMatch ? parseInt(memMatch[1]) : 0;

                this.hostStats.memory = bytes;
                // Convert to GB
                this.hostStats.memoryHuman = (bytes / 1073741824).toFixed(1) + ' GB';

            } catch (e) {
                console.error("Failed to fetch host stats", e);
            }
        },

        async scanLimits() {
            this.showToast('Scanning Resource Limits...');

            // Initialize
            this.jails.forEach(j => {
                if(!this.resourcesData[j.Name]) this.resourcesData[j.Name] = { limits: [], quota: 'Loading...', loading: true };
                else this.resourcesData[j.Name].loading = true;
            });

            // 1. Fetch RCTL Limits
            this.queueTask(async () => {
                try {
                    const endpoint = `${this.apiConfig.url}/api/v1/bastille/limits?target=ALL&action=list`;
                    const response = await this.fetchWithAuth(endpoint, 'POST');
                    const text = await response.text();

                    // Parse: [jail]: \n ... \n jail:name:resource=value:action
                    const lines = text.trim().split('\n');
                    let currentJail = null;
                    const limitMap = {};

                    lines.forEach(line => {
                        line = line.trim();
                        if (!line || line.startsWith('---') || line.includes('[RCTL Limits]')) return;

                        const jailMatch = line.match(/^\[([a-zA-Z0-9\-\._]+)\]:$/);
                        if (jailMatch) {
                            currentJail = jailMatch[1];
                            if(!limitMap[currentJail]) limitMap[currentJail] = [];
                            return;
                        }

                        // RCTL Line: jail:proxy:memoryuse=1G:deny
                        if (currentJail && line.startsWith('jail:')) {
                            const parts = line.split(':');
                            // Expected parts: [0]jail, [1]name, [2]resource=value, [3]action
                            if (parts.length >= 4) {
                                // Parse resource=value
                                const [resource, value] = parts[2].split('=');
                                limitMap[currentJail].push({
                                    resource: resource,
                                    value: value,
                                    action: parts[3],
                                    raw: line
                                });
                            }
                        }
                    });

                    // Assign limits to state
                    Object.keys(this.resourcesData).forEach(name => {
                        this.resourcesData[name].limits = limitMap[name] || [];
                    });

                } catch (e) {
                    console.error("Limit scan failed", e);
                }
            });

            // 2. Fetch ZFS Quotas (Queue per jail)
            this.jails.forEach(jail => {
                this.queueTask(async () => {
                    try {
                        const endpoint = `${this.apiConfig.url}/api/v1/bastille/zfs?target=${jail.Name}&action=get&key_value=quota`;
                        const response = await fetch(endpoint, { method: 'POST', headers: { 'Authorization': `Bearer ${this.apiConfig.key}`, 'Authorization-ID': `bastille` } });
                        if(response.ok) {
                            const text = await response.text();
                            // ZFS output: "NAME PROPERTY VALUE SOURCE" -> "zroot/.../root refquota 1G local"
                            // We grab the value (3rd column usually, or just look for the size pattern/none)
                            // Regex: match "none" or digit+unit at the end of line, or as a standalone token
                            const match = text.match(/\b(none|[\d.]+[BKMGTPE]?)\b/i);
                            const val = match ? match[0] : 'Unknown';
                            this.resourcesData[jail.Name].quota = (val === 'none' || val === '-') ? 'NO LIMIT' : val;

                            // Map "Used" from existing storage data
                            // Looks for dataset ending in "/jailName" (e.g. zroot/bastille/jails/proxy)
                            const storageItem = this.storage.find(s => s.name.endsWith('/' + jail.Name));
                            this.resourcesData[jail.Name].used = storageItem ? storageItem.used : '0B';
                        }
                    } catch(e) {
                        this.resourcesData[jail.Name].quota = 'Error';
                    } finally {
                        this.resourcesData[jail.Name].loading = false;
                    }
                });
            });
        },

        async applyLimit(jailName) {
            if (!this.newLimit.value) return;

            this.showToast(`Applying ${this.newLimit.resource} limit to ${jailName}...`);

            // Bastille syntax: bastille limits TARGET set resource value
            // Note: The API likely maps 'action=set' to the 'set' subcommand.
            // We need to construct the args. Based on typical bastille usage: "memoryuse 1G"
            try {
                const args = `${this.newLimit.resource} ${this.newLimit.value}`;
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/limits?target=${jailName}&action=set&args=${encodeURIComponent(args)}`;

                await this.fetchWithAuth(endpoint, 'POST');
                this.showToast("Limit applied.");
                setTimeout(() => this.scanLimits(), 1000);
            } catch (e) {
                this.showToast(`Error: ${e.message}`);
            }
        },

        async openImportModal() {
            // Reset Form
            this.newImport = {
                file: '',
                release: this.releases.length > 0 ? this.releases[0] : '', // Default to latest release
                force: false,
                staticMac: false,
                mac: ''
            };

            this.showImportModal = true;
            this.backups = []; // Clear old list

            // Fetch available backups
            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/list?item=backups`;
                const response = await this.fetchWithAuth(endpoint, 'POST');
                const text = await response.text();
                this.backups = text.trim().split('\n').filter(line => line.length > 0);

                if (this.backups.length > 0) {
                    this.newImport.file = this.backups[0];
                }
            } catch (e) {
                this.showToast(`Failed to load backups: ${e.message}`);
            }
        },

        async importJail() {
            if (!this.newImport.file) return;

            this.showToast(`Importing ${this.newImport.file}...`);
            this.showImportModal = false;

            try {
                // Build Options: -f -M
                // Strip patch version from release (e.g. 14.3-RELEASE-p1 -> 14.3-RELEASE)
                const cleanRelease = this.newImport.release.replace(/-p\d+$/, '');

                let args = '';
                if (this.newImport.force) args += ' -f';
                if (this.newImport.staticMac) {
                    args += ' -M';
                    if (this.newImport.mac) args += ' ' + this.newImport.mac;
                }

                const endpoint = `${this.apiConfig.url}/api/v1/bastille/import?file=${this.newImport.file}&release=${cleanRelease}&arg=${args.trim()}`;
                await this.fetchWithAuth(endpoint, 'POST');

                this.showToast('Import successful.');
                await this.refreshData();
            } catch (e) {
                this.showToast(`Import Failed: ${e.message}`);
            }
        },

        getTemplateIcon(templateName) {
            const category = templateName.split('/')[0];
            const map = {
                'databases': 'fa-solid fa-database',
                'devel': 'fa-solid fa-code',
                'dns': 'fa-solid fa-sitemap',
                'games': 'fa-solid fa-gamepad',
                'java': 'fa-brands fa-java',
                'lang': 'fa-solid fa-terminal',
                'mail': 'fa-solid fa-envelope',
                'misc': 'fa-solid fa-cubes',
                'multimedia': 'fa-solid fa-photo-film',
                'net-mgmt': 'fa-solid fa-list-check',
                'net-p2p': 'fa-solid fa-share-nodes',
                'net': 'fa-solid fa-network-wired',
                'security': 'fa-solid fa-shield-halved',
                'sysutils': 'fa-solid fa-screwdriver-wrench',
                'textproc': 'fa-solid fa-file-lines',
                'www': 'fa-solid fa-globe'
            };

            // Check specific category map, otherwise check if it looks like a custom git repo (often http/git/ssh)
            // If strictly following the "category/name" format of ports, anything else defaults to custom.
            return map[category] || 'fa-solid fa-cube';
        },

        async manageNetwork() {
            const { target, action, interface: iface, ip, vlan, flags } = this.networkForm;

            if (!target || !iface) {
                this.showToast('Target Jail and Interface are required.');
                return;
            }

            // Build Options String
            let opts = '';
            if (flags.auto) opts += ' -a';

            // Add-specific flags
            if (action === 'add') {
                if (flags.bridge) opts += ' -B';
                if (flags.staticMac) opts += ' -M';
                if (flags.noIp) opts += ' -n';
                if (flags.passthrough) opts += ' -P';
                if (flags.vnet) opts += ' -V';
                if (vlan) opts += ` -v ${vlan}`;
            }

            this.showToast(`Performing network ${action} on ${target}...`);

            try {
                // API: /api/v1/bastille/network?target=X&action=add&iface=Y&ip=Z&options=...
                const params = new URLSearchParams();
                params.append('target', target);
                params.append('action', action);
                params.append('iface', iface);
                if (ip && action === 'add') params.append('ip', ip);
                if (opts) params.append('options', opts.trim());

                const endpoint = `${this.apiConfig.url}/api/v1/bastille/network?${params.toString()}`;
                await this.fetchWithAuth(endpoint, 'POST');

                this.showToast(`Network command sent successfully.`);
                // Refresh to potentially see IP changes (though iface add/remove might not reflect immediately in list without re-scan)
                setTimeout(() => this.refreshData(), 1500);
            } catch (e) {
                this.showToast(`Network Action Failed: ${e.message}`);
            }
        },

        async setDiskQuota(jailName, value) {
            this.showToast(`Setting quota for ${jailName}...`);
            try {
                const endpoint = `${this.apiConfig.url}/api/v1/bastille/zfs?target=${jailName}&action=set&key_value=quota=${value}`;
                await this.fetchWithAuth(endpoint, 'POST');
                this.showToast("Quota updated.");
                setTimeout(() => this.scanLimits(), 1000); // Re-scan to update UI
            } catch (e) {
                this.showToast(`Error: ${e.message}`);
            }
        },

        async generateInsights(targetJail = null) {
            const targets = targetJail ? [targetJail] : this.jails.map(j => j.Name);

            targets.forEach(name => {
                if (!this.insightsData[name]) this.insightsData[name] = { ps: '', socks: '', sysrc: '', last: '', mount: '', cron: '' };
                this.insightsData[name].loading = true;
            });

            targets.forEach(name => {
                this.queueTask(async () => {
                    try {
                    // 1. Fetch Processes (ps -aux)
                    const urlPs = `${this.apiConfig.url}/api/v1/bastille/cmd?target=${name}&command=ps -aux`;
                    const resPs = await this.fetchWithAuth(urlPs, 'POST');
                    const txtPs = await resPs.text();

                    // 2. Fetch Sockets (sockstat -4)
                    const urlSock = `${this.apiConfig.url}/api/v1/bastille/cmd?target=${name}&command=sockstat -4`;
                    const resSock = await this.fetchWithAuth(urlSock, 'POST');
                    const txtSock = await resSock.text();

                    // 3. Fetch Sysrc (sysrc -ae)
                    const urlSysrc = `${this.apiConfig.url}/api/v1/bastille/sysrc?target=${name}&args=-ae`;
                    const resSysrc = await this.fetchWithAuth(urlSysrc, 'POST');
                    const txtSysrc = await resSysrc.text();

                    // 4. Fetch Last Login (last)
                    const urlLast = `${this.apiConfig.url}/api/v1/bastille/cmd?target=${name}&command=last`;
                    const resLast = await this.fetchWithAuth(urlLast, 'POST');
                    const txtLast = await resLast.text();

                    // 5. Fetch Mounts (mount)
                    const urlMount = `${this.apiConfig.url}/api/v1/bastille/cmd?target=${name}&command=mount`;
                    const resMount = await this.fetchWithAuth(urlMount, 'POST');
                    const txtMount = await resMount.text();

                    // 6. Fetch Cron (crontab -l)
                    const urlCron = `${this.apiConfig.url}/api/v1/bastille/cmd?target=${name}&command=crontab -l`;
                    const cronHeaders = {};
                    if (this.apiConfig.key) cronHeaders['Authorization'] = `Bearer ${this.apiConfig.key}`;
                    headers['Authorization-ID'] = 'bastille';

                    // Use raw fetch because crontab -l returns exit code 1 (HTTP 500) if the table is empty
                    const resCron = await fetch(urlCron, { method: 'POST', headers: cronHeaders });
                    const txtCron = resCron.ok ? await resCron.text() : 'No root crontab configured.';


                    this.insightsData[name] = {
                        loading: false,
                        ps: txtPs,
                        socks: txtSock,
                        sysrc: txtSysrc,
                        last: txtLast,
                        mount: txtMount,
                        cron: txtCron
                    };
                } catch (e) {
                        this.insightsData[name] = { loading: false, ps: `Error: ${e.message}`, socks: '', sysrc: '', last: '', mount: '', cron: '' };
                    }
                });
            });
        },

        parseSecurityOutput(text) {
            // Look for "0 problem(s)" or similar standard pkg output
            const isClean = text.includes('0 problem(s)');
            // Rough count of "is vulnerable:" occurrences
            const vulnMatches = (text.match(/is vulnerable:/g) || []).length;
            
            if (isClean) return { status: 'secure', count: 0 };
            if (vulnMatches > 0) return { status: 'vulnerable', count: vulnMatches };
            return { status: 'unknown', count: 0 };
        },

        async refreshData() {
            this.loading = true;

            if (this.apiConfig.url) {
                try {
                    const headers = {};
                    if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
                    headers['Authorization-ID'] = `bastille`;

                    // Parallel requests
                    const pJails = fetch(`${this.apiConfig.url}/api/v1/bastille/list?options=-j`, { method: 'POST', headers });
                    const pReleases = fetch(`${this.apiConfig.url}/api/v1/bastille/list?item=releases`, { method: 'POST', headers });
                    const pTemplates = fetch(`${this.apiConfig.url}/api/v1/bastille/list?item=templates`, { method: 'POST', headers });
                    const pStorage = fetch(`${this.apiConfig.url}/api/v1/bastille/zfs?target=ALL&action=usage`, { method: 'POST', headers: headers });
                    
                    // FIX: 'args' instead of 'param' for config command to avoid 400 error
                    const pNetwork = fetch(`${this.apiConfig.url}/api/v1/bastille/config?target=ALL&action=get&property=ip4.addr`, { method: 'POST', headers });

                    const [rJails, rReleases, rTemplates, rStorage, rNetwork] = await Promise.all([pJails, pReleases, pTemplates, pStorage, pNetwork]);

                    // --- 1. Jails (JSON) ---
                    if (rJails.ok) {
                        const data = await rJails.json();
                        if (Array.isArray(data)) {
                            this.jails = data.map((jail) => ({
                                JID: jail.JID,
                                Name: jail.Name,
                                Boot: jail.Boot,
                                State: jail.State,
                                Type: jail.Type,
                                Prio: jail.Prio,
                                IP: jail["IP Address"],
                                Release: jail.Release,
                                Ports: jail["Published Ports"] !== '-' ? jail["Published Ports"] : '',
                                Tags: jail.Tags !== '-' ? jail.Tags : ''
                            }));
                        } else {
                            console.warn("API returned non-array for jails list:", data);
                            this.jails = [];
                        }
                    }

                    // --- 2. Releases (Text) ---
                    if (rReleases.ok) {
                        const text = await rReleases.text();
                        this.releases = text.trim().split('\n')
                            .filter(line => line.length > 0)
                            // Sort descending (Newest/Largest number first)
                            .sort((a, b) => {
                                return parseFloat(b) - parseFloat(a);
                            });
                    }

                    // --- 3. Templates (Text) ---
                    if (rTemplates.ok) {
                        const text = await rTemplates.text();
                        const rawList = text.trim().split('\n').filter(line => line.length > 0 && line.includes('/') && !line.startsWith('/'));
                        this.templates = rawList.map(name => ({
                            name: name,
                            source: 'Remote Registry',
                            description: 'Template fetched from API'
                        })).sort((a, b) => a.name.localeCompare(b.name));
                    }

                    // --- 4. Storage (Text/ZFS) ---
                    if (rStorage.ok) {
                        const text = await rStorage.text();
                        this.storage = this.parseStorageOutput(text);

                        // Find the dataset with the LARGEST available space to represent Host Storage
                        let maxBytes = 0;
                        let maxStr = '-';
                        this.storage.forEach(ds => {
                            const bytes = this.parseBytes(ds.avail);
                            if (bytes > maxBytes) { maxBytes = bytes; maxStr = ds.avail; }
                        });
                        this.hostStats.storage = maxStr;
                    }

                    // --- 5. Network Map ---
                    if (rNetwork.ok) {
                        const text = await rNetwork.text();
                        this.buildNetworkMap(text);
                    }

                    // Fetch Host Stats lazily if on resources page, or just once
                    if (this.currentView === 'resources' || !this.hostStats.ncpu) {
                        this.fetchHostStats();
                    }

                    this.showToast('Data refreshed from API');

                } catch (error) {
                    console.error("Fetch failed:", error);
                    this.showToast(`Connection Failed: ${error.message}`);
                    this.apiConnected = false;
                }
            } else {
                 // Mock data fallback
                 await new Promise(r => setTimeout(r, 600)); 
            }
            this.loading = false;
        },

        // Helper to DRY up fetch calls
        async fetchWithAuth(url, method = 'POST') {
            const headers = {};
            if (this.apiConfig.key) headers['Authorization'] = `Bearer ${this.apiConfig.key}`;
            headers['Authorization-ID'] = `bastille`;
            const response = await fetch(url, { method: method, headers: headers });
            if (!response.ok) throw new Error(`API Error: ${response.status}`);
            return response;
        },

        showToast(msg) {
            this.toast.message = msg;
            this.toast.visible = true;
            setTimeout(() => { this.toast.visible = false }, 3000);
        }
    }
}
