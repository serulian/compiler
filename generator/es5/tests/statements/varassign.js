$module('varassign', function () {
  var $static = this;
  $static.DoSomething = function () {
    var i;
    var $state = $t.sm(function ($continue) {
      i = $t.box(2, $g.____testlib.basictypes.Integer);
      i = $t.box(3, $g.____testlib.basictypes.Integer);
      $t.box(1234, $g.____testlib.basictypes.Integer);
      $state.resolve();
    });
    return $promise.build($state);
  };
});
