$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    return $g.basic.someBool;
  };
  $static.TEST = function () {
    return $g.basic.anotherBool;
  };
  this.$init(function () {
    return $promise.new(function (resolve) {
      $static.someBool = $t.fastbox(true, $g.________testlib.basictypes.Boolean);
      resolve();
    });
  }, 'af4b3683', []);
  this.$init(function () {
    return $promise.new(function (resolve) {
      $static.anotherBool = $g.basic.AnotherFunction();
      resolve();
    });
  }, 'c4254aac', ['af4b3683']);
});
