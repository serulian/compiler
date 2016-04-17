$module('varnoinit', function () {
  var $static = this;
  $static.DoSomething = function () {
    var i;
    var $state = $t.sm(function ($continue) {
      $t.box(1234, $g.____testlib.basictypes.Integer);
      $state.resolve();
    });
    return $promise.build($state);
  };
});
