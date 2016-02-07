$module('boolean', function () {
  var $static = this;
  $static.TEST = function () {
    var first;
    var second;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            first = $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
            second = $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
            $state.resolve($t.nominalwrap((($t.nominalunwrap(first) && $t.nominalunwrap(second)) || $t.nominalunwrap(first)) || !$t.nominalunwrap(second), $g.____testlib.basictypes.Boolean));
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
