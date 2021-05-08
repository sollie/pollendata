const puppeteer = require('puppeteer');

puppeteer.launch({
  executablePath: '/usr/bin/chromium-browser',
  headless: true,
  args: [
    '--no-sandbox', '--disable-gpu'
  ]
}).then(async browser => {
  const page = await browser.newPage();
  let url = 'https://www.naaf.no/pollenvarsel/';
  await page.goto(url, {'waitUntil': 'networkidle0', 'timeout': 0}).catch(e => {
    console.log(e.toString());
  });

  var pollenNames = await page.evaluate(() => pollenNames);
  var area = await page.evaluate(() => areas);

  var today = await page.evaluate(() => document.querySelector('#todaysdate').innerText);
  var tomorrow = await page.evaluate(() => document.querySelector('#tomorrowsdate').innerText);

  await browser.close();

  var data = {areas: area, names: pollenNames, dates: {today: today.slice(-10), tomorrow: tomorrow.slice(-10)}};

  console.log(JSON.stringify(data));

});
