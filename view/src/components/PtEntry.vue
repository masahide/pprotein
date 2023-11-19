<script lang="ts">
import { defineComponent } from "vue";
// import TsvTable from "./TsvTable.vue";

export default defineComponent({
  data() {
    return {
      data: "",
    };
  },
  async beforeCreate() {
    const resp = await fetch(`/api/pt/${this.$route.params.id}`);
    this.data = await resp.text();
  },
});
</script>

<template>
  <!-- コードブロック -->
  <div class="pt-entry">
    <pre>{{ data }}</pre>
  </div>
</template>