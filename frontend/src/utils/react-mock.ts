// Mock React for sonner library (Vue-only app)
// This prevents sonner from trying to import React at runtime
export default {}
export const createElement = () => null
export const Fragment = () => null
export const useState = () => [null, () => {}]
export const useEffect = () => {}
export const useRef = () => ({ current: null })
export const forwardRef = (component: any) => component
export const memo = (component: any) => component

