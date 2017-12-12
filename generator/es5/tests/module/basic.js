$module('basic', function () {
  var $static = this;
  $static.AnotherFunction = function () {
    return $g.basic.someBool;
  };
  $static.TEST = function () {
    return $t.fastbox($g.basic.anotherBool.$wrapped && ($g.basic.anotherInt.$wrapped == 42), $g.________testlib.basictypes.Boolean);
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
  }, '0893c862', ['af4b3683']);
  this.$init(function () {
    return $promise.new(function (resolve) {
      $static.anotherInt = $t.fastbox(42, $g.________testlib.basictypes.Integer);
      resolve();
    });
  }, '02c976cb', []);
});
