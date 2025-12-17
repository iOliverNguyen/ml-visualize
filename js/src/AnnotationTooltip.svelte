<script lang="ts">
  interface Props {
    text: string
    visible?: boolean
    position?: 'top' | 'bottom' | 'left' | 'right'
  }

  let { text, visible = true, position = 'top' }: Props = $props()
</script>

{#if visible}
  <div class="annotation-tooltip" class:top={position === 'top'} class:bottom={position === 'bottom'} class:left={position === 'left'} class:right={position === 'right'}>
    <div class="tooltip-content">
      {@html text}
    </div>
    <div class="tooltip-arrow"></div>
  </div>
{/if}

<style>
  .annotation-tooltip {
    position: absolute;
    z-index: 1000;
    pointer-events: none;
    animation: fadeIn 0.2s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(-4px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .tooltip-content {
    background: #1e293b;
    color: #f1f5f9;
    padding: 0.75rem 1rem;
    border-radius: 6px;
    font-size: 0.875rem;
    line-height: 1.5;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    max-width: 300px;
    white-space: normal;
  }

  .tooltip-content :global(strong) {
    color: #60a5fa;
    font-weight: 600;
  }

  .tooltip-content :global(code) {
    background: rgba(255, 255, 255, 0.1);
    padding: 0.125rem 0.25rem;
    border-radius: 3px;
    font-family: 'SF Mono', monospace;
    font-size: 0.85em;
  }

  .tooltip-arrow {
    position: absolute;
    width: 0;
    height: 0;
  }

  /* Position: Top (tooltip above target) */
  .annotation-tooltip.top {
    transform: translate(-50%, -100%) translateY(-8px);
  }

  .annotation-tooltip.top .tooltip-arrow {
    bottom: -6px;
    left: 50%;
    transform: translateX(-50%);
    border-left: 6px solid transparent;
    border-right: 6px solid transparent;
    border-top: 6px solid #1e293b;
  }

  /* Position: Bottom (tooltip below target) */
  .annotation-tooltip.bottom {
    transform: translate(-50%, 0%) translateY(8px);
  }

  .annotation-tooltip.bottom .tooltip-arrow {
    top: -6px;
    left: 50%;
    transform: translateX(-50%);
    border-left: 6px solid transparent;
    border-right: 6px solid transparent;
    border-bottom: 6px solid #1e293b;
  }

  /* Position: Left (tooltip to the left of target) */
  .annotation-tooltip.left {
    transform: translate(-100%, -50%) translateX(-8px);
  }

  .annotation-tooltip.left .tooltip-arrow {
    right: -6px;
    top: 50%;
    transform: translateY(-50%);
    border-top: 6px solid transparent;
    border-bottom: 6px solid transparent;
    border-left: 6px solid #1e293b;
  }

  /* Position: Right (tooltip to the right of target) */
  .annotation-tooltip.right {
    transform: translate(0%, -50%) translateX(8px);
  }

  .annotation-tooltip.right .tooltip-arrow {
    left: -6px;
    top: 50%;
    transform: translateY(-50%);
    border-top: 6px solid transparent;
    border-bottom: 6px solid transparent;
    border-right: 6px solid #1e293b;
  }
</style>
