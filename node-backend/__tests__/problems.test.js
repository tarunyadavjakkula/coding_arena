const request = require('supertest');
const app = require('../src/app');

describe('GET /api/problems/:id', () => {
  test('200 — valid ID returns problem', async () => {
    const res = await request(app).get('/api/problems/1');
    expect(res.status).toBe(200);
    expect(res.body.id).toBe('1');
  });

  test('404 — unknown ID', async () => {
    const res = await request(app).get('/api/problems/999');
    expect(res.status).toBe(404);
  });

  test('400 — non-numeric ID', async () => {
    const res = await request(app).get('/api/problems/abc');
    expect(res.status).toBe(400);
  });
});
