const Market = artifacts.require('Market');
const Point = artifacts.require('Point');

module.exports = async function (deployer) {
  try {
    // await deployer.deploy(Point, 'SUN');
    // await deployer.deploy(Market);
  } catch (error) {
    console.log(error);
  }
  return;
};
