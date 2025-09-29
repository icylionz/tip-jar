// Tip Jar Application JavaScript
// Alpine.js, HTMX, and Tailwind CSS are loaded from CDN in templates

// Application initialization
document.addEventListener("DOMContentLoaded", function () {
  // Configure HTMX if available
  if (typeof htmx !== "undefined") {
    htmx.config.defaultSwapStyle = "outerHTML";
    htmx.config.useTemplateFragments = true;
  }

  // Global application data for Alpine.js
  window.tipjar = {
    user: null,
    currentJar: null,
    notifications: [],
    // Reactive state management
    state: {
      loading: false,
      error: null,
    },
  };
  utils.formatTimestamps();
});

// Utility functions
const utils = {
  formatCurrency: (amount) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(amount);
  },

  formatDate: (date) => {
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(new Date(date));
  },

  formatRelativeTime: (date) => {
    const now = new Date();
    const target = new Date(date);
    const diffInSeconds = Math.floor((now - target) / 1000);

    if (diffInSeconds < 60) return "just now";
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
    if (diffInSeconds < 86400)
      return `${Math.floor(diffInSeconds / 3600)}h ago`;
    if (diffInSeconds < 2592000)
      return `${Math.floor(diffInSeconds / 86400)}d ago`;

    return utils.formatDate(date);
  },
  formatTimestamps: () => {
    document.querySelectorAll("[data-timestamp]").forEach((el) => {
      const timestamp = el.getAttribute("data-timestamp");
      if (timestamp) {
        const formatted = utils.formatRelativeTime(timestamp);
        // Replace the "Joined [date]" part with relative time
        el.textContent = el.textContent.replace(
          /Joined .*/,
          `Joined ${formatted}`
        );
      }
    });
  },
  showNotification: (message, type = "info") => {
    const notification = document.createElement("div");
    notification.className = `notification notification-${type}`;
    notification.innerHTML = `
            <div class="flex items-center">
                <span>${message}</span>
                <button onclick="this.parentElement.parentElement.remove()" class="ml-3 text-lg font-semibold">&times;</button>
            </div>
        `;

    document.body.appendChild(notification);

    // Auto remove after 5 seconds
    setTimeout(() => {
      if (notification.parentElement) {
        notification.remove();
      }
    }, 5000);
  },

  copyToClipboard: async (text) => {
    try {
      await navigator.clipboard.writeText(text);
      utils.showNotification("Copied to clipboard!", "success");
    } catch (err) {
      utils.showNotification("Failed to copy to clipboard", "error");
    }
  },

  generateInviteCode: () => {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    let result = "";
    for (let i = 0; i < 8; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
  },

  validateEmail: (email) => {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
  },

  debounce: (func, wait) => {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  },
};

// Alpine.js components and stores
document.addEventListener("alpine:init", () => {
  // Global notification store
  Alpine.store("notifications", {
    items: [],
    add(message, type = "info") {
      const id = Date.now();
      this.items.push({ id, message, type });
      setTimeout(() => this.remove(id), 5000);
    },
    remove(id) {
      this.items = this.items.filter((item) => item.id !== id);
    },
  });

  // Loading state store
  Alpine.store("loading", {
    state: false,
    show() {
      this.state = true;
    },
    hide() {
      this.state = false;
    },
  });
});

// Export for global use
window.utils = utils;

// Service worker registration for offline support (future enhancement)
if ("serviceWorker" in navigator) {
  window.addEventListener("load", () => {
    // navigator.serviceWorker.register('/static/js/sw.js'); // TODO: Implement service worker
  });
}
// Offense Type Form Data
function offenseTypeFormData(jarID) {
  return {
    form: {
      name: '',
      description: '',
      cost_type: 'monetary',
      cost_amount: '',
      cost_action: ''
    },
    error: null,
    loading: false,
    jarID: jarID,
    
    updateCostType() {
      if (this.form.cost_type === 'monetary') {
        this.form.cost_action = '';
      } else {
        this.form.cost_amount = '';
      }
    },
    
    async submitForm() {
      this.loading = true;
      this.error = null;
      
      try {
        const formData = new FormData();
        formData.append('name', this.form.name);
        formData.append('description', this.form.description);
        formData.append('cost_type', this.form.cost_type);
        
        if (this.form.cost_type === 'monetary') {
          if (!this.form.cost_amount) {
            this.error = 'Cost amount is required for monetary offenses';
            this.loading = false;
            return;
          }
          formData.append('cost_amount', this.form.cost_amount);
        } else {
          if (!this.form.cost_action) {
            this.error = 'Cost action is required';
            this.loading = false;
            return;
          }
          formData.append('cost_action', this.form.cost_action);
        }
        
        const url = `/jars/${this.jarID}/offense-types`;
        console.log('Submitting to URL:', url);
        
        const response = await fetch(url, {
          method: 'POST',
          body: formData
        });
        
        if (response.ok) {
          window.location.reload();
        } else {
          const errorText = await response.text();
          this.error = errorText || 'Failed to save offense type';
        }
      } catch (error) {
        console.error('Caught error:', error);
        this.error = 'Network error. Please try again.';
      } finally {
        this.loading = false;
      }
    }
  };
}