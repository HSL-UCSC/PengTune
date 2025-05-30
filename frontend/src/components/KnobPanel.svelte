<script>
  import { onMount } from "svelte";
  import jQuery from "jquery";
  import "round-slider";

  const axes = ["X", "Y", "Z"];
  const gains = ["P", "I", "D"];
  const groups = ["Position", "Attitude"];

  // All knob definitions: 3 axes × 3 gains × 2 groups = 18 knobs
  const knobDefs = groups.flatMap(group =>
    axes.flatMap(axis =>
      gains.map(gain => ({
        id: `${group.slice(0, 3)}${axis}${gain}`, // e.g., PosXP, AttZD
        label: `${axis} / ${gain}`,
        axis,
        gain,
        group
      }))
    )
  );

  // Knob values
  let values = {};
  knobDefs.forEach(({ id }) => values[id] = 10);

  // Link states: group → gain → true/false (default: true)
  let linkXY = {
    Position: { P: true, I: true, D: true },
    Attitude: { P: true, I: true, D: true }
  };

  function onChange(id, value) {
    values[id] = value;

    // Parse ID
    const group = id.startsWith("Pos") ? "Position" : "Attitude";
    const axis = id[3]; // X, Y, or Z
    const gain = id[4]; // P, I, or D

    // If linking is active and axis is X or Y, sync with the other
    if (linkXY[group][gain] && (axis === "X" || axis === "Y")) {
      const otherAxis = axis === "X" ? "Y" : "X";
      const otherId = `${id.slice(0, 3)}${otherAxis}${gain}`;
      if (values[otherId] !== value) {
        values[otherId] = value;
        jQuery(`#${otherId}`).roundSlider("setValue", value);
        console.log(`Synced ${id} → ${otherId}`);
        window.go.main.App.PublishKnob({ knob: otherId, value });
      }
    }

    // Publish primary knob
    console.log("Knob changed:", id, value);
    window.go.main.App.PublishKnob({ knob: id, value })
      .then(() => console.log("Published", id, value))
      .catch((err) => console.error("PublishKnob failed:", err));
  }

  function updateKnob(id) {
    jQuery(`#${id}`).roundSlider("setValue", values[id]);
    onChange(id, values[id]);
  }

onMount(() => {

  window.runtime.EventsOn("update:pos", (gains) => {
    console.log("update:pos", gains);

    const kp = gains.kp ?? [];
    const ki = gains.ki ?? [];
    const kd = gains.kd ?? [];

    const base = "Pos";
    ["P", "I", "D"].forEach((gain) => {
      const arr = gain === "P" ? kp : gain === "I" ? ki : kd;
      ["X", "Y", "Z"].forEach((axis, j) => {
        const id = `${base}${axis}${gain}`;
        if (arr[j] !== undefined) {
          values[id] = arr[j];
          jQuery(`#${id}`).roundSlider("setValue", arr[j]);
        } else {
          console.warn(`Missing value for ${id}`);
        }
      });
    });
  });

  window.runtime.EventsOn("update:att", (gains) => {
    console.log("update:att", gains);

    const kp = gains.kp ?? [];
    const ki = gains.ki ?? [];
    const kd = gains.kd ?? [];

    const base = "Att";
    ["P", "I", "D"].forEach((gain) => {
      const arr = gain === "P" ? kp : gain === "I" ? ki : kd;
      ["X", "Y", "Z"].forEach((axis, j) => {
        const id = `${base}${axis}${gain}`;
        if (arr[j] !== undefined) {
          values[id] = arr[j];
          jQuery(`#${id}`).roundSlider("setValue", arr[j]);
        } else {
          console.warn(`Missing value for ${id}`);
        }
      });
    });
  });

  knobDefs.forEach(({ id }) => {
    jQuery(`#${id}`).roundSlider({
      radius: 60,
      width: 10,
      min: 0,
      max: 20,
      step: 0.1,
      value: values[id],
      sliderType: "min-range",
      editableTooltip: true,
      showTooltip: true,

      // Clamp during drag to prevent wrap-around
      drag: function(args) {
        const MIN = args.options.min;
        const MAX = args.options.max;
        let v = args.value;

        if (v <= MIN) {
          this.setValue(MIN);
          onChange(id, MIN);
          return false; // stop internal wrap logic
        }
        if (v >= MAX) {
          this.setValue(MAX);
          onChange(id, MAX);
          return false;
        }
        onChange(id, v);
      },

      // Final clamp on release
      change: function(args) {
        const MIN = args.options.min;
        const MAX = args.options.max;
        const v = Math.max(Math.min(args.value, MAX), MIN);
        onChange(id, v);
      }
    });
  });
});
</script>

<style>
  .panel {
    display: flex;
    justify-content: space-around;
    gap: 4rem;
    padding: 2rem;
    flex-wrap: wrap;
  }

  .group {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .group h2 {
    margin-bottom: 0.5rem;
    font-size: 1.25rem;
  }

  .checkboxes {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 1rem;
  }

  .checkboxes label {
    font-size: 0.9rem;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1.25rem;
  }

  .knob-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100px;
  }

  .knob-label {
    margin-top: 0.5rem;
    font-weight: bold;
    font-size: 0.9rem;
    text-align: center;
  }

  .knob-input {
    margin-top: 0.25rem;
    width: 60px;
    text-align: center;
    color: #000
  }

</style>

<div class="panel">
  {#each groups as group}
    <div class="group">
      <h2>{group}</h2>

      <!-- Lock toggles -->
      <div class="checkboxes">
        {#each gains as gain}
          <label>
            <input type="checkbox" bind:checked={linkXY[group][gain]} />
            Link X/Y {gain}
          </label>
        {/each}
      </div>

      <!-- Knob grid -->
      <div class="grid">
        {#each knobDefs.filter(k => k.group === group) as { id, label }}
          <div class="knob-wrapper">
            <div id={id}></div>
            <div class="knob-label">{label}</div>
          </div>
        {/each}
      </div>
    </div>
  {/each}
</div>
