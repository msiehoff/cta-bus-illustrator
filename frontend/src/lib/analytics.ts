const GA_MEASUREMENT_ID = 'G-M6JMBQS6EJ'

declare global {
  interface Window {
    dataLayer?: unknown[]
    gtag?: Gtag
  }
}

type Gtag = {
  (...args: unknown[]): void
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  (command: 'config' | 'event' | 'js' | 'set', ...args: any[]): void
}

let initPromise: Promise<void> | null = null

const setupDataLayer = (): void => {
  window.dataLayer = window.dataLayer || []
  // gtag.js expects `arguments` objects in the dataLayer, not plain arrays
  window.gtag = function gtag() {
    // eslint-disable-next-line prefer-rest-params
    window.dataLayer!.push(arguments)
  } as Gtag
}

export const initAnalytics = (): Promise<void> => {
  if (!import.meta.env.PROD) return Promise.resolve()
  if (initPromise) return initPromise

  setupDataLayer()
  window.gtag!('js', new Date())

  initPromise = new Promise(resolve => {
    const script = document.createElement('script')
    script.async = true
    script.src = `https://www.googletagmanager.com/gtag/js?id=${GA_MEASUREMENT_ID}`
    script.onload = () => resolve()
    script.onerror = () => resolve()
    document.head.appendChild(script)
  })

  return initPromise
}

export const trackPageView = async (pagePath: string): Promise<void> => {
  if (!import.meta.env.PROD) return
  await initAnalytics()
  window.gtag?.('config', GA_MEASUREMENT_ID, { page_path: pagePath })
}
