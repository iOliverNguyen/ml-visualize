import type { Glossary, FAQData } from '../types';

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
 * Loads all educational content at once
 */
export async function loadAllContent() {
  const [glossary, faqs] = await Promise.all([
    loadGlossary(),
    loadFAQs()
  ]);

  return {
    glossary,
    faqs
  };
}
