import test from 'ava';

const Suite = require('./helpers/suite');

const color = require('color');
const palette = require('./helpers/palette');

test.beforeEach(async t => {
  t.context = new Suite();
  await t.context.start(t);
});

test.afterEach(async t => {
  t.context.passed(t);
});

test.always.afterEach(async t => {
  await t.context.finish(t);
});

test('scrolls to the top with gg, and to the bottom with G', async t => {
  await t.context.fly.run('set-pipeline -n -p some-pipeline -c fixtures/pipeline-with-long-output.yml');
  await t.context.fly.run('unpause-pipeline -p some-pipeline');

  await t.context.fly.run('trigger-job -j some-pipeline/long-output');

  await t.context.page.goto(t.context.web.route(`/teams/${t.context.teamName}/pipelines/some-pipeline/jobs/long-output/builds/1`));

  await t.context.page.waitForFunction(() => {
    return document.body.innerText.indexOf("Line 100") !== -1
  }, {
    polling: 100,
    timeout: 60000
  });

  await t.context.page.type('body', 'G');
  await t.context.page.waitForFunction(() => window.scrollY > 0);

  await t.context.page.type('body', 'gg');
  await t.context.page.waitForFunction(() => window.scrollY == 0);

  // need an assertion  *somewhere*
  t.true(true);
});
