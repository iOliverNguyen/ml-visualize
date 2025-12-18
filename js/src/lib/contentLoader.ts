import type { Glossary, FAQData, TutorialContent } from '../types';

/**
 * Loads the glossary from external JSON file
 */
export async function loadGlossary(): Promise<Glossary> {
  try {
    const response = await fetch('/content/glossary.json');
    if (!response.ok) {
      throw new Error(`Failed to load glossary: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error loading glossary:', error);
    return {};
  }
}

/**
 * Loads FAQs from external JSON file
 */
export async function loadFAQs(): Promise<FAQData> {
  try {
    const response = await fetch('/content/faqs.json');
    if (!response.ok) {
      throw new Error(`Failed to load FAQs: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error loading FAQs:', error);
    return { categories: [] };
  }
}

/**
 * Loads the tutorial from external JSON file
 */
export async function loadTutorial(): Promise<TutorialContent> {
  try {
    const response = await fetch('/content/tutorial.json');
    if (!response.ok) {
      throw new Error(`Failed to load tutorial: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error loading tutorial:', error);
    return {
      meta: {
        title: '',
        description: '',
        estimatedReadTime: 0,
        version: '1.0'
      },
      chapters: []
    };
  }
}

/**
 * Loads all educational content at once
 */
export async function loadAllContent() {
  const [glossary, faqs, tutorial] = await Promise.all([
    loadGlossary(),
    loadFAQs(),
    loadTutorial()
  ]);

  return {
    glossary,
    faqs,
    tutorial
  };
}
