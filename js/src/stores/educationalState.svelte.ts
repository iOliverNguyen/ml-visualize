/**
 * Educational state management using Svelte 5 runes
 * Handles user preferences, UI state, and progress tracking
 */

// Internal state (not exported directly due to Svelte 5 restrictions)
let _userLevel = $state<'beginner' | 'intermediate' | 'advanced'>('beginner');
let _showSidebar = $state<boolean>(false);
let _sidebarPosition = $state<'left' | 'right'>('right');
let _expandedSections = $state<Set<string>>(new Set());
let _activeGlossaryTerm = $state<string | null>(null);
let _tourActive = $state<boolean>(false);
let _currentTourStep = $state<number>(0);
let _tourCompleted = $state<boolean>(false);
let _showIntroPanel = $state<boolean>(true);

// Export getters/setters for reactive access
export const state = {
  get userLevel() { return _userLevel; },
  set userLevel(value) { _userLevel = value; },

  get showSidebar() { return _showSidebar; },
  set showSidebar(value) { _showSidebar = value; },

  get sidebarPosition() { return _sidebarPosition; },
  set sidebarPosition(value) { _sidebarPosition = value; },

  get expandedSections() { return _expandedSections; },
  set expandedSections(value) { _expandedSections = value; },

  get activeGlossaryTerm() { return _activeGlossaryTerm; },
  set activeGlossaryTerm(value) { _activeGlossaryTerm = value; },

  get tourActive() { return _tourActive; },
  set tourActive(value) { _tourActive = value; },

  get currentTourStep() { return _currentTourStep; },
  set currentTourStep(value) { _currentTourStep = value; },

  get tourCompleted() { return _tourCompleted; },
  set tourCompleted(value) { _tourCompleted = value; },

  get showIntroPanel() { return _showIntroPanel; },
  set showIntroPanel(value) { _showIntroPanel = value; },
};

/**
 * Toggle a section's expanded state
 */
export function toggleSection(sectionId: string) {
  if (_expandedSections.has(sectionId)) {
    _expandedSections.delete(sectionId);
  } else {
    _expandedSections.add(sectionId);
  }
  // Trigger reactivity
  _expandedSections = _expandedSections;
}

/**
 * Check if a section is expanded
 */
export function isSectionExpanded(sectionId: string): boolean {
  return _expandedSections.has(sectionId);
}

/**
 * Set user level
 */
export function setUserLevel(level: 'beginner' | 'intermediate' | 'advanced') {
  _userLevel = level;
  saveToLocalStorage();
}

/**
 * Toggle sidebar visibility
 */
export function toggleSidebar() {
  _showSidebar = !_showSidebar;
  saveToLocalStorage();
}

/**
 * Save state to localStorage
 */
export function saveToLocalStorage() {
  if (typeof localStorage !== 'undefined') {
    try {
      localStorage.setItem('educational-state', JSON.stringify({
        userLevel: _userLevel,
        showSidebar: _showSidebar,
        sidebarPosition: _sidebarPosition,
        expandedSections: Array.from(_expandedSections),
        tourCompleted: _tourCompleted,
        showIntroPanel: _showIntroPanel
      }));
    } catch (e) {
      console.warn('Failed to save educational state:', e);
    }
  }
}

/**
 * Load state from localStorage
 */
export function loadFromLocalStorage() {
  if (typeof localStorage !== 'undefined') {
    try {
      const saved = localStorage.getItem('educational-state');
      if (saved) {
        const savedState = JSON.parse(saved);
        if (savedState.userLevel) _userLevel = savedState.userLevel;
        if (typeof savedState.showSidebar === 'boolean') _showSidebar = savedState.showSidebar;
        if (savedState.sidebarPosition) _sidebarPosition = savedState.sidebarPosition;
        if (Array.isArray(savedState.expandedSections)) {
          _expandedSections = new Set(savedState.expandedSections);
        }
        if (typeof savedState.tourCompleted === 'boolean') _tourCompleted = savedState.tourCompleted;
        if (typeof savedState.showIntroPanel === 'boolean') _showIntroPanel = savedState.showIntroPanel;
      }
    } catch (e) {
      console.warn('Failed to load educational state:', e);
    }
  }
}

// Auto-load on module initialization
if (typeof window !== 'undefined') {
  loadFromLocalStorage();
}
