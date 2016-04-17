$module('nullcompare', function () {
  var $static = this;
  $static.TEST = function () {
    var someBool;
    var $state = $t.sm(function ($continue) {
      someBool = null;
      $state.resolve($t.nullcompare(someBool, $t.box(true, $g.____testlib.basictypes.Boolean)));
      return;
    });
    return $promise.build($state);
  };
});
