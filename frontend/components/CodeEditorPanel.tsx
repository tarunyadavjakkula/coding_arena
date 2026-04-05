import { useRef, useEffect, useCallback, memo } from 'react'
import type { OnMount } from '@monaco-editor/react'
import type { editor } from 'monaco-editor'
import { lazy, Suspense } from 'react'

// Lazy load Monaco Editor using React lazy
const Editor = lazy(() => import('@monaco-editor/react'))

export interface LanguageTemplate {
  id: string
  name: string
  monacoId: string
  extension: string
  starterCode: string
}

interface CodeEditorPanelProps {
  value: string
  language: string
  onChange: (value: string) => void
  onLanguageChange?: (language: string) => void
  availableLanguages?: LanguageTemplate[]
  readOnly?: boolean
  debounceDelay?: number // Debounce delay in milliseconds (default: 300)
}

function CodeEditorPanelComponent({
  value,
  language,
  onChange,
  onLanguageChange,
  availableLanguages = [],
  readOnly = false,
  debounceDelay = 50, // Default 50ms debounce delay
}: CodeEditorPanelProps) {
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null)
  const debounceTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const handleEditorDidMount: OnMount = (editor) => {
    editorRef.current = editor
  }

  // Debounced onChange handler to reduce re-renders
  const handleEditorChange = useCallback((value: string | undefined) => {
    // Clear existing timer
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current)
    }
    
    // Set new timer to debounce onChange events
    debounceTimerRef.current = setTimeout(() => {
      onChange(value || '')
    }, debounceDelay)
  }, [onChange, debounceDelay])

  const handleLanguageSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newLanguage = e.target.value
    if (onLanguageChange) {
      onLanguageChange(newLanguage)
    }
  }

  // Cleanup on unmount to prevent memory leaks
  useEffect(() => {
    return () => {
      // Clear debounce timer
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current)
      }
      // Dispose editor safely
      if (editorRef.current) {
        try {
          editorRef.current.dispose()
          editorRef.current = null
        } catch (error) {
          console.error('Error disposing Monaco editor:', error)
        }
      }
    }
  }, [])

  return (
    <div className="h-full w-full flex flex-col bg-white">
      {availableLanguages.length > 0 && onLanguageChange && (
        <div className="bg-gray-100 border-b border-gray-200 px-4 py-2 flex items-center gap-2">
          <label htmlFor="language-selector" className="text-gray-700 text-sm font-medium">
            Language:
          </label>
          <select
            id="language-selector"
            value={language}
            onChange={handleLanguageSelect}
            className="bg-white text-gray-900 border border-gray-300 rounded px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-green-500 transition-all duration-200"
          >
            {availableLanguages.map((lang) => (
              <option key={lang.id} value={lang.id}>
                {lang.name}
              </option>
            ))}
          </select>
        </div>
      )}
      <div className="flex-1 min-h-0 relative">
        <Suspense fallback={
          <div className="h-full w-full flex items-center justify-center bg-white text-gray-800">
            <div className="text-center">
              <div className="flex justify-center mb-4">
                <div className="w-12 h-12 border-4 border-green-500 border-t-transparent rounded-full animate-spin"></div>
              </div>
              <div className="text-lg mb-2">Loading Editor...</div>
            </div>
          </div>
        }>
          <Editor
            height="100%"
            language={language}
            value={value}
            theme="light"
            onChange={handleEditorChange}
            onMount={handleEditorDidMount}
            loading={
              <div className="h-full w-full flex items-center justify-center bg-white text-gray-800">
                <div className="text-center">
                  <div className="flex justify-center mb-4">
                    <div className="w-12 h-12 border-4 border-green-500 border-t-transparent rounded-full animate-spin"></div>
                  </div>
                  <div className="text-lg mb-2">Loading Monaco Editor...</div>
                  <div className="text-gray-600 text-sm">Please wait</div>
                </div>
              </div>
            }
            options={{
              fontSize: 14,
              lineNumbers: 'on',
              minimap: { enabled: true },
              automaticLayout: true,
              scrollBeyondLastLine: false,
              wordWrap: 'off',
              tabSize: 4,
              readOnly,
            }}
          />
        </Suspense>
      </div>
    </div>
  )
}

// Memoize component to prevent unnecessary re-renders
export const CodeEditorPanel = memo(CodeEditorPanelComponent)
