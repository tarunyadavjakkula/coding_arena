import { useMemo, memo } from 'react'
import DOMPurify from 'isomorphic-dompurify'

export interface ConsoleMessage {
  type: 'info' | 'success' | 'error'
  message: string
  timestamp: number
}

interface ConsolePanelProps {
  messages: ConsoleMessage[]
  onClose?: () => void
  status?: string
  submissionId?: string
}

function ConsolePanelComponent({ messages, onClose, status, submissionId }: ConsolePanelProps) {
  const sanitizedMessages = useMemo(
    () =>
      messages.map((msg) => ({
        ...msg,
        message: DOMPurify.sanitize(msg.message),
      })),
    [messages]
  )

  const getMessageStyle = (type: ConsoleMessage['type']) => {
    switch (type) {
      case 'info':
        return 'text-blue-600 bg-blue-50 border-blue-200'
      case 'success':
        return 'text-green-600 bg-green-50 border-green-200'
      case 'error':
        return 'text-red-600 bg-red-50 border-red-200'
      default:
        return 'text-gray-600 bg-gray-50 border-gray-200'
    }
  }

  return (
    <div className="h-full w-full flex flex-col bg-white">
      <div className="bg-gray-100 border-b border-gray-200 px-4 py-2 flex items-center justify-between">
        <div className="flex items-center gap-4 flex-1">
          <h2 className="text-gray-900 text-sm font-semibold">Console</h2>
          {status && (
            <span className="text-gray-900 text-sm font-semibold">{status}</span>
          )}
          {submissionId && (
            <span className="text-gray-900 text-sm font-semibold">{submissionId}</span>
          )}
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className="text-gray-600 hover:text-gray-900 transition-colors duration-200 p-1 rounded hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-green-500"
            title="Close console"
            aria-label="Close console"
          >
            <span className="text-lg leading-none">×</span>
          </button>
        )}
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-3">
        {sanitizedMessages.length === 0 ? (
          <div className="text-gray-500 text-sm">
            Console output will appear here
          </div>
        ) : (
          sanitizedMessages.map((msg, index) => (
            <div
              key={index}
              className={`rounded border px-4 py-3 text-sm ${getMessageStyle(msg.type)}`}
            >
              <div
                dangerouslySetInnerHTML={{ __html: msg.message }}
              />
            </div>
          ))
        )}
      </div>
    </div>
  )
}

// Memoize component to prevent unnecessary re-renders
export const ConsolePanel = memo(ConsolePanelComponent)
