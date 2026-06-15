const GA_MEASUREMENT_ID = 'G-M6JMBQS6EJ'

declare global {
  interface Window {
    dataLayer?: unknown[]
    gtag?: (...args: unknown[]) => void
  }
}

let initialized = false

export const initAnalytics = (): void => {
  if (!import.meta.env.PROD || initialized) return
  initialized = true

  window.dataLayer = window.dataLayer || []
  window.gtag = (...args: unknown[]) => {
    window.dataLayer!.push(args)
  }
  window.gtag('js', new Date())
  window.gtag('config', GA_MEASUREMENT_ID)

  const script = document.createElement('script')
  script.async = true
  script.src = `https://www.googletagmanager.com/gtag/js?id=${GA_MEASUREMENT_ID}`
  document.head.appendChild(script)
}

export const trackPageView = (pagePath: string): void => {
  if (!import.meta.env.PROD) return
  window.gtag?.('config', GA_MEASUREMENT_ID, { page_path: pagePath })
}
