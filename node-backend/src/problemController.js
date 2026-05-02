const { getProblemById } = require('./problemService');

const getProblem = (req, res) => {
  try {
    const { id } = req.params;

    if (!id || !/^\d+$/.test(id)) {
      return res.status(400).json({ error: 'Invalid problem ID provided. ID must be numeric.' });
    }

    if (!Number.isSafeInteger(Number(id)) || Number(id) <= 0) {
      return res.status(400).json({ error: 'Invalid problem ID provided. ID must be numeric.' });
    }

    const problem = getProblemById(id);

    if (!problem) {
      return res.status(404).json({ error: `Problem with ID ${id} not found` });
    }

    return res.status(200).json(problem);

  } catch (error) {
    console.error(`Error fetching problem ${req.params.id}:`, error);
    return res.status(500).json({ error: 'An internal server error occurred' });
  }
};

module.exports = {
  getProblem
};
