// TODO: Replace this in-memory mock with real YAML file reads from
// judge-config/problems/<id>.yaml using the 'js-yaml' package.
// This ensures the Node backend and Go judge use the same problem definitions.
const problems = [
  {
    id: "1",
    title: "Pair Index Match",
    difficulty: "Easy",
    statement: "Given a list of integers and a goal value, find the positions of two different elements whose sum equals the goal.",
    constraints: [
      "2 <= values.length <= 10^4",
      "-10^9 <= values[i] <= 10^9",
      "-10^9 <= goal <= 10^9",
      "Assume exactly one matching pair is present."
    ],
    timeLimit: 1000, // milliseconds
    memoryLimit: 256 // megabytes
  },
  {
    id: "2",
    title: "Linked List Digit Sum",
    difficulty: "Medium",
    statement: "Two singly linked lists each encode a non-negative integer using one digit per node in least-significant-first order. Produce a new linked list representing the total in the same format.",
    constraints: [
      "Each input list contains between 1 and 100 nodes.",
      "Each node stores a single digit from 0 to 9.",
      "The encoded values may be zero or larger and should be combined without changing the input lists."
    ],
    timeLimit: 2000,
    memoryLimit: 256
  }
];

module.exports = problems;
