$module('conditional', function () {
  var $static = this;
  $static.DoSomething = function () {
    return $t.fastbox(42, $g.____testlib.basictypes.Integer);
  };
  $static.TEST = function () {
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          if ($g.conditional.DoSomething().$wrapped == 42) {
            $current = 1;
            continue syncloop;
          } else {
            $current = 2;
            continue syncloop;
          }
          break;

        case 1:
          return $t.fastbox(true, $g.____testlib.basictypes.Boolean);

        case 2:
          return $t.fastbox(false, $g.____testlib.basictypes.Boolean);

        default:
          return;
      }
    }
  };
});
