# ADR-006: Use Vue 3 with Static Site Generation

## Status

Accepted

## Date

2025-10-21

## Context

The application requires a modern frontend that:

- Displays a booking calendar with available slots
- Allows users to request bookings
- Provides admin interface for managing bookings
- Communicates with GraphQL API
- Authenticates via AWS Cognito
- Is deployed to S3 + CloudFront for minimal hosting costs

Frontend requirements:

- Responsive design for mobile and desktop
- Fast loading times
- Modern developer experience
- Easy styling with utility CSS
- Search engine optimization is NOT a priority (private booking links)

The development team has experience with Vue.js and prefers to continue using it. The team wants to use Tailwind CSS for rapid UI development.

## Decision

We will use **Vue 3** with **Static Site Generation (SSG)** via **Vite** and **Tailwind CSS** for styling.

Implementation approach:

- Vue 3 Composition API for component logic
- Vite for build tooling and development server
- Vue Router for client-side routing
- Pinia for state management
- Tailwind CSS for styling
- Static build output deployed to S3
- Client-side GraphQL calls to API Gateway
- AWS Amplify libraries for Cognito authentication

## Consequences

### Positive

- **Zero Hosting Cost**: Static files on S3 with CloudFront
  - No server-side rendering compute costs
  - S3 costs: ~$0.023 per GB/month (likely <$0.10/month)
  - CloudFront: 1TB free tier data transfer
  - **Expected cost: ~$0/month**
- **Fast Loading**: Pre-built static assets
  - HTML/CSS/JS served instantly from CDN
  - No server-side rendering delays
  - CloudFront edge caching provides <100ms response times
- **Simple Deployment**: Just upload files to S3
  - No server configuration
  - No container orchestration
  - Easy CI/CD integration
- **Team Familiarity**: Team has Vue experience
  - Faster development velocity
  - Less learning curve
  - Comfortable with Vue ecosystem
- **Modern Development Experience**:
  - Vite provides instant hot module replacement (HMR)
  - Fast build times (seconds vs. minutes)
  - Great TypeScript support
- **Vue 3 Benefits**:
  - Composition API for better code organization
  - Better TypeScript integration
  - Smaller bundle sizes than Vue 2
  - Improved performance
- **Tailwind CSS Benefits**:
  - Rapid UI development with utility classes
  - Consistent design system
  - Automatic purging of unused CSS (tiny bundle sizes)
  - Responsive design made easy
  - No CSS naming conflicts
- **Client-Side Rendering Appropriate**: SEO not important
  - Application is accessed via direct links
  - No need for search engine indexing
  - Simplifies architecture significantly

### Negative

- **No True SSR/SSG**: Despite "SSG" name, this is actually client-side rendering
  - Initial HTML is empty shell
  - JavaScript required for app to function
  - Mitigation: Acceptable since SEO not required; users have modern browsers
- **Initial Load Time**: Must download Vue + app bundles before rendering
  - Typical bundle size: 100-200KB gzipped
  - On 3G network: 1-2 seconds
  - Mitigation: Low traffic users likely have good connections; CloudFront caching helps
- **No Progressive Enhancement**: JavaScript required for functionality
  - App won't work with JS disabled
  - Mitigation: Acceptable trade-off for modern web app
- **Client-Side API Calls**: All GraphQL requests from browser
  - Exposes API endpoint to users
  - Must handle authentication tokens in browser
  - Mitigation: JWT tokens in localStorage/sessionStorage; API Gateway validates all requests
- **Build Step Required**: Cannot edit and deploy directly
  - Must run build process for changes
  - Mitigation: Automated via CI/CD; fast build times with Vite
- **Tailwind Learning Curve**: Team must learn utility-first CSS
  - Different mindset from traditional CSS
  - Can lead to very long class lists
  - Mitigation: Team eager to learn; good documentation available

### Technology Stack Details

**Core Framework**:

- Vue 3.4+ (Composition API)
- Vite 5+ (build tool and dev server)
- TypeScript (for type safety)

**Routing & State**:

- Vue Router 4 (client-side routing)
- Pinia (state management store)

**Styling**:

- Tailwind CSS 3+ (utility-first CSS)
- PostCSS (for Tailwind processing)
- Autoprefixer (browser compatibility)

**UI Components**:

- Start with custom components
- Consider Headless UI or Radix Vue for accessibility later

**GraphQL Client**:

- Apollo Client for Vue or urql
- Automatic caching and request management
- TypeScript types generated from schema

**Authentication**:

- AWS Amplify Auth library
- Handles Cognito integration
- Token management and refresh

**Build Output**:

```text
dist/
  index.html
  assets/
    index-[hash].js
    index-[hash].css
  favicon.ico
```

### Development Workflow

1. **Local Development**:

   ```bash
   npm run dev  # Vite dev server on localhost:5173
   ```

2. **Build**:

   ```bash
   npm run build  # Creates dist/ folder
   ```

3. **Preview**:

   ```bash
   npm run preview  # Test production build locally
   ```

4. **Deploy**:
   - Upload dist/ to S3 via Terraform/GitHub Actions
   - CloudFront invalidation for cache clearing

### Project Structure

```text
frontend/
  src/
    assets/          # Images, fonts
    components/      # Vue components
      BookingCalendar.vue
      BookingForm.vue
      AdminDashboard.vue
    views/           # Page components
      HomePage.vue
      BookingsPage.vue
      AdminPage.vue
    router/          # Vue Router config
      index.ts
    store/           # Pinia stores
      auth.ts
      bookings.ts
    graphql/         # GraphQL queries/mutations
      queries.ts
      mutations.ts
    utils/           # Utilities
      auth.ts
      date.ts
    App.vue
    main.ts
  public/            # Static assets (copied as-is)
  index.html         # HTML template
  tailwind.config.js
  vite.config.ts
  package.json
```

### Example Component

```vue
<template>
  <div class="max-w-4xl mx-auto p-6">
    <h1 class="text-3xl font-bold text-gray-900 mb-6">
      Book Caselo di Salzan
    </h1>

    <BookingCalendar
      v-model:selected-date="selectedDate"
      :available-slots="availableSlots"
      @date-click="handleDateClick"
    />

    <BookingForm
      v-if="selectedDate"
      :date="selectedDate"
      @submit="handleBookingSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { GET_AVAILABILITY } from '@/graphql/queries'
import BookingCalendar from '@/components/BookingCalendar.vue'
import BookingForm from '@/components/BookingForm.vue'

const selectedDate = ref<string | null>(null)

const { result, loading } = useQuery(GET_AVAILABILITY, {
  start: '2025-11-01',
  end: '2025-11-30'
})

const availableSlots = computed(() => result.value?.availability ?? [])

const handleDateClick = (date: string) => {
  selectedDate.value = date
}

const handleBookingSubmit = async (bookingData: any) => {
  // Submit booking mutation
}
</script>
```

## Alternatives Considered

### Nuxt 3 (Vue SSR Framework)

- **Pros**:
  - True SSR/SSG capabilities
  - Better SEO (not needed here)
  - Faster initial page load
  - File-based routing
  - Built-in API routes
- **Cons**:
  - More complex setup
  - Requires server or Edge runtime for SSR
  - Higher hosting costs (not truly static if using SSR)
  - SSG still requires build-time data fetching (not suitable for dynamic bookings)
  - Overkill for client-side app
- **Rejected**: Additional complexity not justified; no SEO requirement; SSR costs money

### React + Next.js

- **Pros**:
  - Larger ecosystem
  - More job market demand
  - Excellent SSR/SSG support
- **Cons**:
  - Team unfamiliar with React
  - Steeper learning curve combined with Go backend learning
  - Similar SSR/SSG cost issues as Nuxt
- **Rejected**: Team preference for Vue; no added value for this use case

### Svelte + SvelteKit

- **Pros**:
  - Smaller bundle sizes
  - No virtual DOM (better performance)
  - Simple syntax
- **Cons**:
  - Team unfamiliar with Svelte
  - Smaller ecosystem
  - Less mature than Vue
  - Additional learning curve
- **Rejected**: Team preference for Vue; learning Go already is enough new tech

### Plain HTML/CSS/JavaScript

- **Pros**:
  - No framework overhead
  - Maximum performance
  - No build step
- **Cons**:
  - Manually manage DOM updates
  - No component reusability
  - No state management
  - Complex to maintain
  - Slow development
- **Rejected**: Too primitive for modern app with dynamic calendar and forms

### jQuery + Bootstrap

- **Pros**:
  - Simple and familiar to many developers
  - Quick prototyping
- **Cons**:
  - Outdated approach
  - Poor developer experience
  - Imperative DOM manipulation
  - Not suitable for modern SPAs
- **Rejected**: Not modern enough; Vue provides much better DX

## Migration Path

If client-side rendering proves problematic (unlikely):

1. **Add Nuxt 3** for SSR/SSG
2. **Reuse Vue components** (minimal changes needed)
3. **Deploy to Edge runtime** (Cloudflare Workers, Vercel Edge)
4. **Accept higher hosting costs** (~$5-20/month)

Estimated migration time: 1 week

## SEO Considerations

While SEO is not a priority, basic considerations:

- **Meta Tags**: Add proper title and description in index.html
- **Social Sharing**: Open Graph tags for link previews
- **robots.txt**: Allow/disallow as needed
- **sitemap.xml**: Can be generated and placed in public/

If SEO becomes important later, migrate to Nuxt 3 SSR/SSG.

## Performance Budget

Target metrics:

- First Contentful Paint (FCP): <1.5s
- Largest Contentful Paint (LCP): <2.5s
- Time to Interactive (TTI): <3.5s
- Total Bundle Size: <200KB gzipped

Strategies:

- Code splitting by route
- Lazy load components
- Tailwind CSS purging (removes unused styles)
- Image optimization
- CloudFront compression and caching

## References

- [Vue 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
- [Pinia State Management](https://pinia.vuejs.org/)
- [AWS Amplify for Vue](https://docs.amplify.aws/vue/)
- [Vue Apollo GraphQL](https://apollo.vuejs.org/)
