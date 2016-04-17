$module('assertnotnull', function () {
  var $static = this;
  $static.TEST = function () {
    var someValue;
    var $state = $t.sm(function ($continue) {
      someValue = $t.box(true, $g.____testlib.basictypes.Boolean);
      $state.resolve($t.assertnotnull(someValue));
      return;
    });
    return $promise.build($state);
  };
});
