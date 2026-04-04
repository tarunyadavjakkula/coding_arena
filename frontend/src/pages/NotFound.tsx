import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div className="h-screen w-screen flex items-center justify-center bg-white text-gray-800 p-4">
      <div className="text-center max-w-md">
        <h1 className="text-5xl sm:text-6xl font-semibold mb-4 text-gray-900">404</h1>
        <h2 className="text-xl sm:text-2xl mb-2 font-semibold text-gray-900">Problem Not Found</h2>
        <p className="text-gray-600 mb-6 text-sm sm:text-base">
          The problem you're looking for doesn't exist.
        </p>
        <Link
          to="/problems"
          className="inline-block px-6 py-3 bg-gray-900 hover:bg-gray-800 text-white rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-gray-400"
        >
          Back to Problems
        </Link>
      </div>
    </div>
  )
}
