$module('init', function () {
  var $static = this;
  $static.TEST = function () {
    return $t.fastbox(($g.init.sc.value.$wrapped == 2) && ($g.init.sc2.value.$wrapped == 4), $g.________testlib.basictypes.Boolean);
  };
  this.$init(function () {
    return $promise.new(function (resolve) {
      $static.sc = $g.other.SomeClass.NewThing($t.fastbox(1, $g.________testlib.basictypes.Integer));
      resolve();
    });
  }, '194711ba', ['cc19450d']);
  this.$init(function () {
    return $promise.new(function (resolve) {
      $static.sc2 = $g.other.SomeClass.NewThing($t.fastbox(2, $g.________testlib.basictypes.Integer));
      resolve();
    });
  }, '2e80db22', ['cc19450d']);
});
