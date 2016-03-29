$module('assertnotnull', function () {
  var $static = this;
  $static.TEST = function () {
    var someValue;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            someValue = $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
            $state.resolve($t.assertnotnull(someValue));
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
