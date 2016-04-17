$module('boolean', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($continue) {
      $t.box(true, $g.____testlib.basictypes.Boolean);
      $t.box(false, $g.____testlib.basictypes.Boolean);
      $state.resolve();
    });
    return $promise.build($state);
  };
});
