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

test('shows abort hooks', async t => {
  await t.context.fly.run('set-pipeline -n -p some-pipeline -c fixtures/hooks-pipeline.yml');
  await t.context.fly.run('unpause-pipeline -p some-pipeline');

  await t.context.fly.run('trigger-job -j some-pipeline/on_abort');

  await t.context.page.goto(t.context.web.route(`/teams/${t.context.teamName}/pipelines/some-pipeline/jobs/on_abort/builds/1`));
  await t.context.web.waitForText(t.context.page, "say-bye-from-step");
  await t.context.web.waitForText(t.context.page, "say-bye-from-job");
  await t.context.web.waitForText(t.context.page, "looping");

  await t.context.web.clickAndWait(t.context.page, '.build-action-abort');
  await t.context.page.waitFor('[data-step-name="say-bye-from-step"] i.succeeded');
  await t.context.page.waitFor('[data-step-name="say-bye-from-job"] i.succeeded');

  await t.context.web.clickAndWait(t.context.page, '[data-step-name="say-bye-from-step"] .header');
  t.regex(await t.context.web.text(t.context.page), /bye from step/);

  await t.context.web.clickAndWait(t.context.page, '[data-step-name="say-bye-from-job"] .header');
  t.regex(await t.context.web.text(t.context.page), /bye from job/);
});

test('can be switched between', async t => {
  await t.context.fly.run('set-pipeline -n -p some-pipeline -c fixtures/states-pipeline.yml');
  await t.context.fly.run('unpause-pipeline -p some-pipeline');

  await t.context.fly.run('trigger-job -w -j some-pipeline/passing');
  await t.context.fly.run('trigger-job -w -j some-pipeline/passing');

  await t.context.page.goto(t.context.web.route(`/teams/${t.context.teamName}/pipelines/some-pipeline/jobs/passing/builds/1`));

  await t.context.web.clickAndWait(t.context.page, '#builds li:nth-child(1) a');
  t.regex(await t.context.web.text(t.context.page), /passing #2/);

  await t.context.web.clickAndWait(t.context.page, '#builds li:nth-child(2) a');
  t.regex(await t.context.web.text(t.context.page), /passing #1/);
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
