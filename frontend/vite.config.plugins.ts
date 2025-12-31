import type { Plugin } from 'vite'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

/**
 * Plugin para mockar React no runtime
 * Necessário porque sonner tenta importar React dinamicamente
 */
export function reactMockPlugin(): Plugin {
  return {
    name: 'react-mock',
    enforce: 'pre',
    resolveId(id) {
      // Interceptar importações de React
      if (id === 'react' || id === 'react-dom') {
        const mockPath = path.resolve(
          path.dirname(fileURLToPath(import.meta.url)),
          'src/utils/react-mock.ts'
        )
        return mockPath
      }
      return null
    },
    load(id) {
      // Se for o arquivo mock, retornar o conteúdo
      if (id.includes('react-mock.ts')) {
        return `
// Mock React for sonner library (Vue-only app)
const createElement = (type, props, ...children) => {
  // Return a simple object that won't break sonner
  return { type, props: props || {}, children }
}

const Fragment = ({ children }) => children

export default {
  createElement,
  Fragment,
  useState: () => [null, () => {}],
  useEffect: () => {},
  useRef: () => ({ current: null }),
  forwardRef: (component) => component,
  memo: (component) => component,
  createContext: () => ({ Provider: () => null, Consumer: () => null }),
  useContext: () => null,
  useMemo: (fn) => fn(),
  useCallback: (fn) => fn,
  useReducer: () => [null, () => {}],
  useLayoutEffect: () => {},
  useImperativeHandle: () => {},
  useDebugValue: () => {},
  Component: class {},
  PureComponent: class {},
  StrictMode: () => null,
  Suspense: () => null,
  lazy: (fn) => fn,
  createPortal: () => null,
  version: '18.0.0'
}

export { createElement, Fragment }
export const useState = () => [null, () => {}]
export const useEffect = () => {}
export const useRef = () => ({ current: null })
export const forwardRef = (component) => component
export const memo = (component) => component
export const createContext = () => ({ Provider: () => null, Consumer: () => null })
export const useContext = () => null
export const useMemo = (fn) => fn()
export const useCallback = (fn) => fn
export const useReducer = () => [null, () => {}]
export const useLayoutEffect = () => {}
export const useImperativeHandle = () => {}
export const useDebugValue = () => {}
export const Component = class {}
export const PureComponent = class {}
export const StrictMode = () => null
export const Suspense = () => null
export const lazy = (fn) => fn
export const createPortal = () => null
export const version = '18.0.0'
        `
      }
      return null
    },
  }
}

