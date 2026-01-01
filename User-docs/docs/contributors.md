---
title: Contributors
---

# Contributors

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Contributor {
  login: string
  avatar_url: string
  html_url: string
  name?: string | null
}

const contributors = ref<Contributor[]>([])
const error = ref<string | null>(null)
const loading = ref(true)

async function fetchContributors() {
  try {
    const res = await fetch(
      'https://api.github.com/repos/Naganathan05/Load-Pulse/contributors?per_page=100'
    )

    if (!res.ok) {
      throw new Error('Failed to load contributors from GitHub.')
    }

    const data = await res.json()

    // Fetch full user profiles to get display names
    const detailed = await Promise.all(
      (data as any[]).map(async (c) => {
        try {
          const userRes = await fetch(`https://api.github.com/users/${c.login}`)
          if (!userRes.ok) return c
          const user = await userRes.json()
          return { ...c, name: user.name || c.login }
        } catch {
          return c
        }
      })
    )

    contributors.value = detailed as Contributor[]
  } catch (e: any) {
    error.value = e?.message ?? 'Unable to load contributors.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchContributors()
})
</script>

<div v-if="loading">
  Loading contributors...
</div>

<div v-else-if="error">
  {{ error }}
</div>

<div v-else class="contributors-grid">
  <a
    v-for="c in contributors"
    :key="c.login"
    class="contributor-card"
    :href="c.html_url"
    target="_blank"
    rel="noopener noreferrer"
  >
    <img
      class="contributor-avatar"
      :src="c.avatar_url"
      :alt="c.login"
      loading="lazy"
    />
    <div class="contributor-name">
      {{ c.name || c.login }}
    </div>
    <div class="contributor-username">
      @{{ c.login }}
    </div>
  </a>
</div>

<style scoped>
.contributors-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 1.5rem;
  margin-top: 1.5rem;
}

.contributor-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  border-radius: 0.75rem;
  border: 1px solid var(--vp-c-divider);
  background-color: var(--vp-c-bg-soft);
  text-decoration: none;
  color: inherit;
  transition: transform 0.15s ease, box-shadow 0.15s ease,
    border-color 0.15s ease;
}

.contributor-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  border-color: var(--vp-c-brand);
}

.contributor-avatar {
  width: 80px;
  height: 80px;
  border-radius: 999px;
  margin-bottom: 0.75rem;
  border: 2px solid var(--vp-c-brand);
}

.contributor-name {
  font-weight: 600;
  margin-bottom: 0.25rem;
  text-align: center;
}

.contributor-username {
  font-size: 0.85rem;
  color: var(--vp-c-text-2);
}
</style>
