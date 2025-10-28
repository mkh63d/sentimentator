<script setup>
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <header class="flex items-center max-h-screen leading-normal pr-[calc(var(--section-gap)/2)]">
    <div class="flex items-start flex-wrap w-full">
      <nav
        class="w-full text-xs text-center mt-8 md:text-left md:-ml-4 md:text-base md:py-4 md:mt-4"
      >
        <RouterLink
          to="/"
          class="inline-block px-4 border-l first:border-0 border-(--color-border) hover:bg-transparent [&.router-link-exact-active]:text-(--color-text)"
        >
          Home
        </RouterLink>
        <RouterLink
          :to="{ name: 'login' }"
          class="inline-block px-4 border-l first:border-0 border-(--color-border) hover:bg-transparent [&.router-link-exact-active]:text-(--color-text)"
        >
          Login
        </RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
</template>
