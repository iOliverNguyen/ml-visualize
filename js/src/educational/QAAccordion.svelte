<script lang="ts">
  import ExpandableSection from './ExpandableSection.svelte';
  import type { FAQCategory } from '../types';
  import * as educationalState from '../stores/educationalState.svelte';

  interface Props {
    categories: FAQCategory[];
    filterByLevel?: boolean;
    searchQuery?: string;
  }

  let {
    categories,
    filterByLevel = true,
    searchQuery = ''
  }: Props = $props();

  // Filter questions based on user level and search query
  const filteredCategories = $derived.by(() => {
    const userLevel = educationalState.state.userLevel;
    const query = searchQuery.toLowerCase();

    return (categories ?? []).map(category => {
      let questions = category?.questions ?? [];

      // Filter by level
      if (filterByLevel) {
        questions = questions.filter(q => {
          if (userLevel === 'beginner') {
            return q?.level === 'beginner';
          } else if (userLevel === 'intermediate') {
            return q?.level === 'beginner' || q?.level === 'intermediate';
          } else {
            return true; // Advanced sees all
          }
        });
      }

      // Filter by search query
      if (query) {
        questions = questions.filter(q =>
          q?.question?.toLowerCase().includes(query) ||
          q?.answer?.toLowerCase().includes(query) ||
          (q?.tags ?? []).some(tag => tag?.toLowerCase().includes(query))
        );
      }

      return {
        ...category,
        questions
      };
    }).filter(category => (category?.questions?.length ?? 0) > 0);
  });

  function handleExpandChange(questionId: string, expanded: boolean) {
    if (expanded) {
      educationalState.state.expandedSections.add(`faq-${questionId}`);
    } else {
      educationalState.state.expandedSections.delete(`faq-${questionId}`);
    }
    educationalState.state.expandedSections = educationalState.state.expandedSections;
    educationalState.saveToLocalStorage();
  }
</script>

<div class="qa-accordion">
  {#if filteredCategories().length === 0}
    <div class="empty-state">
      <p>No questions match your current filters.</p>
      {#if searchQuery}
        <p class="hint">Try a different search term.</p>
      {/if}
    </div>
  {:else}
    {#each filteredCategories() as category}
      <div class="category-section">
        <h3 class="category-title">{category.name}</h3>
        <div class="questions">
          {#each category.questions as question}
            <ExpandableSection
              id={`faq-${question.id}`}
              title={question.question}
              defaultExpanded={educationalState.isSectionExpanded(`faq-${question.id}`)}
              level={question.level as 'beginner' | 'intermediate' | 'advanced'}
              icon="❓"
              onexpandChange={(expanded) => handleExpandChange(question.id, expanded)}
            >
              {#snippet children()}
                <div class="answer">
                  <p>{question.answer}</p>
                  {#if question.tags.length > 0}
                    <div class="tags">
                      {#each question.tags as tag}
                        <span class="tag">{tag}</span>
                      {/each}
                    </div>
                  {/if}
                </div>
              {/snippet}
            </ExpandableSection>
          {/each}
        </div>
      </div>
    {/each}
  {/if}
</div>

<style>
  .qa-accordion {
    width: 100%;
  }

  .category-section {
    margin-bottom: 2rem;
  }

  .category-title {
    margin: 0 0 1rem 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: #0f172a;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #e2e8f0;
  }

  .questions {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .answer {
    line-height: 1.6;
  }

  .answer p {
    margin: 0 0 1rem 0;
    color: #334155;
  }

  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid #e2e8f0;
  }

  .tag {
    display: inline-flex;
    align-items: center;
    padding: 0.25rem 0.5rem;
    font-size: 0.75rem;
    font-weight: 500;
    color: #64748b;
    background: #f1f5f9;
    border-radius: 4px;
    text-transform: lowercase;
  }

  .empty-state {
    padding: 2rem;
    text-align: center;
    color: #64748b;
  }

  .empty-state p {
    margin: 0.5rem 0;
  }

  .empty-state .hint {
    font-size: 0.9rem;
    color: #94a3b8;
  }
</style>
