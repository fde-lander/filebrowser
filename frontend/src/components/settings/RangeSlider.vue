<template>
  <div class="range-slider-container">
    <div class="range-slider-label">
      <span class="range-slider-name">{{ name }}</span>
      <span v-if="description" class="range-slider-description">{{ description }}</span>
    </div>
    <div class="range-slider-control">
      <input
        type="range"
        class="range-slider-input"
        :min="min"
        :max="max"
        :step="step"
        :value="modelValue"
        @input="updateValue"
      />
      <span class="range-slider-value">{{ formattedValue }}</span>
    </div>
  </div>
</template>

<script>
export default {
  name: "RangeSlider",
  props: {
    modelValue: {
      type: Number,
      required: true,
    },
    min: {
      type: Number,
      default: 0,
    },
    max: {
      type: Number,
      default: 100,
    },
    step: {
      type: Number,
      default: 1,
    },
    displayFormat: {
      type: String,
      default: "number", // "number" or "percent"
    },
    name: {
      type: String,
      default: "",
    },
    description: {
      type: String,
      default: "",
    },
  },
  emits: ["update:modelValue"],
  computed: {
    formattedValue() {
      if (this.displayFormat === "percent") {
        return `${Math.round(this.modelValue * 100)}%`;
      }
      return String(this.modelValue);
    },
  },
  methods: {
    updateValue(event) {
      const val = parseFloat(event.target.value);
      this.$emit("update:modelValue", val);
    },
  },
};
</script>

<style scoped>
.range-slider-container {
  display: flex;
  flex-direction: column;
  padding: 0.5em 0;
  font-size: 1rem;
}

.range-slider-label {
  display: flex;
  align-items: center;
  margin-bottom: 0.3em;
}

.range-slider-name {
  font-weight: 500;
}

.range-slider-description {
  margin-left: 0.5em;
  font-size: 0.85em;
  opacity: 0.7;
}

.range-slider-control {
  display: flex;
  align-items: center;
  gap: 0.75em;
}

.range-slider-input {
  flex: 1;
  cursor: pointer;
  accent-color: var(--primaryColor);
}

.range-slider-value {
  min-width: 3em;
  text-align: right;
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}
</style>
