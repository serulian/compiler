$module('isnull', function () {
  var $static = this;
  $static.DoSomething = function (a) {
    var $state = $t.sm(function ($continue) {
      $state.resolve($t.box(a == null, $g.____testlib.basictypes.Boolean));
      return;
    });
    return $promise.build($state);
  };
});
