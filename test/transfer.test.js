'use strict';

const test = require('node:test');
const assert = require('node:assert');
const { createLedger } = require('../src/transfer');

test('transfer moves balance between accounts', () => {
  const ledger = createLedger({ alice: 100, bob: 0 });
  ledger.transfer('alice', 'bob', 30);
  assert.equal(ledger.balanceOf('alice'), 70);
  assert.equal(ledger.balanceOf('bob'), 30);
});

test('balanceOf defaults unknown accounts to zero', () => {
  const ledger = createLedger();
  assert.equal(ledger.balanceOf('nobody'), 0);
});
