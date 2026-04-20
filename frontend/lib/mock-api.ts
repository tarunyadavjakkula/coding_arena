/**
 * Mock API interceptor for Vite
 * This simulates backend API responses by intercepting fetch calls.
 *
 * DISABLED by default — the Vite proxy forwards /api/* to the real backend.
 * To re-enable for offline development, set VITE_USE_MOCK=true in your .env
 * or start vite with: VITE_USE_MOCK=true npm run dev
 */

interface RunRequestBody {
  problem_id?: string;
  language?: string;
  source?: string;
  custom_input?: string | null;
}

const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

if (typeof window !== 'undefined' && USE_MOCK) {
  console.warn('[mock-api] Mock API interceptor is ACTIVE. API calls will NOT reach the real backend.')

  const originalFetch = window.fetch
  window.fetch = async (...args) => {
    const [url, options] = args
    const urlString = typeof url === 'string' ? url : url instanceof URL ? url.toString() : ''

    if (urlString.startsWith('/api/submit') && (!options || options.method === 'POST')) {
      console.log('[mock-api] Intercepting /api/submit')

      try {
        await new Promise(resolve => setTimeout(resolve, 800))

        const cases = [
          { position: 1, status: 'AC', time: 0.012, memory_kb: 3200, points: 5.0, total_points: 5.0 },
          { position: 2, status: 'AC', time: 0.015, memory_kb: 3200, points: 5.0, total_points: 5.0 },
          { position: 3, status: 'WA', time: 0.011, memory_kb: 3100, points: 0.0, total_points: 5.0 },
          { position: 4, status: 'AC', time: 0.013, memory_kb: 3200, points: 5.0, total_points: 5.0 },
          { position: 5, status: 'TLE', time: 2.001, memory_kb: 4500, points: 0.0, total_points: 5.0 },
        ]

        const points = cases.reduce((s, c) => s + c.points, 0)
        const totalPoints = cases.reduce((s, c) => s + c.total_points, 0)
        const allAC = cases.every(c => c.status === 'AC')

        const responseData = {
          id: 'sub_' + Math.random().toString(16).slice(2, 34).padEnd(32, '0'),
          status: allAC ? 'graded' : 'graded',
          problem_id: 'mock-problem',
          language: 'python',
          message: allAC ? 'AC' : 'WA',
          result: {
            verdict: allAC ? 'AC' : 'WA',
            cases: cases,
            total_time: cases.reduce((s, c) => s + c.time, 0),
            max_memory_kb: Math.max(...cases.map(c => c.memory_kb)),
            points: points,
            total_points: totalPoints,
          }
        }

        return new Response(JSON.stringify(responseData), {
          status: 200,
          headers: { 'Content-Type': 'application/json' }
        })
      } catch {
        return new Response(JSON.stringify({ error: 'Invalid request' }), {
          status: 400,
          headers: { 'Content-Type': 'application/json' }
        })
      }
    }

    if (urlString.startsWith('/api/run') && (!options || options.method === 'POST')) {
      console.log('[mock-api] Intercepting /api/run')

      try {
        await new Promise(resolve => setTimeout(resolve, 500))

        let requestBody: RunRequestBody = {}
        if (options && options.body) {
          requestBody = JSON.parse(options.body as string)
        }

        const isCustomInput = requestBody.custom_input && requestBody.custom_input.trim()

        let testCases
        if (isCustomInput) {
          testCases = [
            {
              name: 'Custom Input',
              status: 'AC',
              time: 0.010,
              memory_kb: 3200,
              input: requestBody.custom_input,
              expected_output: 'N/A',
              actual_output: 'Your output here (simulated)'
            }
          ]
        } else {
          testCases = [
            {
              name: 'Test Case 1',
              status: 'AC',
              time: 0.012,
              memory_kb: 3200,
              input: '',
              expected_output: '',
              actual_output: ''
            },
            {
              name: 'Test Case 2',
              status: 'AC',
              time: 0.015,
              memory_kb: 3100,
              input: '',
              expected_output: '',
              actual_output: ''
            }
          ]
        }

        const responseData = {
          run_id: 'run_' + Math.random().toString(16).slice(2, 34).padEnd(32, '0'),
          status: 'AC',
          message: isCustomInput ? 'Custom input executed' : 'AC',
          test_cases: testCases,
          timestamp: 0
        }

        return new Response(JSON.stringify(responseData), {
          status: 200,
          headers: { 'Content-Type': 'application/json' }
        })
      } catch {
        return new Response(JSON.stringify({ error: 'Invalid request' }), {
          status: 400,
          headers: { 'Content-Type': 'application/json' }
        })
      }
    }

    return originalFetch(...args)
  }
}

export {}
