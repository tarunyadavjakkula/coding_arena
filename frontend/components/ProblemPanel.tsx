import { useMemo, memo } from 'react'
import { useNavigate } from 'react-router-dom'
import DOMPurify from 'isomorphic-dompurify'
import type { Problem } from '@/lib/data-service'

interface ProblemPanelProps {
  problem: Problem
}

function ProblemPanelComponent({ problem }: ProblemPanelProps) {
  const navigate = useNavigate()

  // Sanitize HTML content to prevent XSS attacks
  const sanitizedDescription = useMemo(
    () => DOMPurify.sanitize(problem.description),
    [problem.description]
  )

  const sanitizedInputFormat = useMemo(
    () => DOMPurify.sanitize(problem.inputFormat),
    [problem.inputFormat]
  )

  const sanitizedOutputFormat = useMemo(
    () => DOMPurify.sanitize(problem.outputFormat),
    [problem.outputFormat]
  )

  return (
    <div className="h-full overflow-y-auto bg-white text-gray-900 p-4 sm:p-6">
      {/* Back Button */}
      <button 
        onClick={() => navigate('/problems')}
        className="mb-4 text-gray-700 hover:text-gray-900 flex items-center gap-2 transition-colors text-sm bg-white hover:bg-gray-50 py-2 px-4 rounded-md border border-gray-300 font-medium"
      >
        <span>←</span> Back to Problems
      </button>

      {/* Problem Title */}
      <h1 className="text-xl sm:text-2xl font-bold mb-4 text-gray-900">{problem.title}</h1>

      {/* Problem Metadata */}
      <div className="flex items-center gap-3 sm:gap-4 mb-6 text-sm flex-wrap">
        <span
          className={`font-semibold ${
            problem.difficulty === 'Easy'
              ? 'text-green-600'
              : problem.difficulty === 'Medium'
              ? 'text-yellow-600'
              : 'text-red-600'
          }`}
        >
          {problem.difficulty}
        </span>
        <span className="text-gray-600">{problem.category}</span>
        <span className="text-gray-600">{problem.points} points</span>
      </div>

      {/* Problem Description */}
      <section className="mb-6">
        <h2 className="text-base sm:text-lg font-bold mb-3 text-gray-900">Description</h2>
        <div
          className="text-gray-700 leading-relaxed text-sm sm:text-base"
          dangerouslySetInnerHTML={{ __html: sanitizedDescription }}
        />
      </section>

      {/* Input Format */}
      <section className="mb-6">
        <h2 className="text-base sm:text-lg font-bold mb-3 text-gray-900">Input Format</h2>
        <div
          className="text-gray-700 leading-relaxed text-sm sm:text-base"
          dangerouslySetInnerHTML={{ __html: sanitizedInputFormat }}
        />
      </section>

      {/* Output Format */}
      <section className="mb-6">
        <h2 className="text-base sm:text-lg font-bold mb-3 text-gray-900">Output Format</h2>
        <div
          className="text-gray-700 leading-relaxed text-sm sm:text-base"
          dangerouslySetInnerHTML={{ __html: sanitizedOutputFormat }}
        />
      </section>

      {/* Constraints */}
      <section className="mb-6">
        <h2 className="text-base sm:text-lg font-bold mb-3 text-gray-900">Constraints</h2>
        <ul className="list-disc list-inside text-gray-700 space-y-1 text-sm sm:text-base">
          {problem.constraints.map((constraint, index) => (
            <li key={index}>{constraint}</li>
          ))}
        </ul>
      </section>

      {/* Examples */}
      <section className="mb-6">
        <h2 className="text-base sm:text-lg font-bold mb-3 text-gray-900">Examples</h2>
        <div className="space-y-4">
          {problem.examples.map((example, index) => (
            <div
              key={index}
              className="bg-gray-50 rounded-lg p-3 sm:p-4 border border-gray-200"
            >
              <h3 className="text-sm font-bold text-gray-700 mb-2">
                Example {index + 1}
              </h3>

              {/* Input */}
              <div className="mb-3">
                <div className="text-xs font-bold text-gray-600 mb-1">
                  Input:
                </div>
                <pre className="bg-white rounded p-2 text-xs sm:text-sm text-gray-800 overflow-x-auto border border-gray-200 font-mono">
                  {example.input}
                </pre>
              </div>

              {/* Output */}
              <div className="mb-3">
                <div className="text-xs font-bold text-gray-600 mb-1">
                  Output:
                </div>
                <pre className="bg-white rounded p-2 text-xs sm:text-sm text-gray-800 overflow-x-auto border border-gray-200 font-mono">
                  {example.output}
                </pre>
              </div>

              {/* Explanation (optional) */}
              {example.explanation && (
                <div>
                  <div className="text-xs font-bold text-gray-600 mb-1">
                    Explanation:
                  </div>
                  <p className="text-xs sm:text-sm text-gray-700 leading-relaxed">{example.explanation}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      </section>

      {/* Time and Memory Limits */}
      <section className="mb-6">
        <h2 className="text-base sm:text-lg font-bold mb-3 text-gray-900">Limits</h2>
        <div className="flex gap-4 sm:gap-6 text-sm text-gray-700 flex-wrap">
          <div>
            <span className="text-gray-600">Time Limit:</span>{' '}
            <span className="font-semibold">{problem.timeLimit}ms</span>
          </div>
          <div>
            <span className="text-gray-600">Memory Limit:</span>{' '}
            <span className="font-semibold">{problem.memoryLimit}MB</span>
          </div>
        </div>
      </section>
    </div>
  )
}

// Memoize component to prevent unnecessary re-renders
export const ProblemPanel = memo(ProblemPanelComponent)
