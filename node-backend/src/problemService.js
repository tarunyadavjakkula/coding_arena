const problems = require('./mockData');

/**
 * Helper function to fetch a problem by its ID
 * @param {string} id - The ID of the problem to fetch
 * @returns {Object|null} - The problem object or null if not found
 */
const getProblemById = (id) => {
  const problem = problems.find(p => p.id === String(id));
  return problem || null;
};

module.exports = {
  getProblemById
};
