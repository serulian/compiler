$module('loopexpr', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          $t.fastbox(1234, $g.____testlib.basictypes.Integer);
          $current = 1;
          continue syncloop;

        case 1:
          if (true) {
            $current = 2;
            continue syncloop;
          } else {
            $current = 3;
            continue syncloop;
          }
          break;

        case 2:
          $t.fastbox(1357, $g.____testlib.basictypes.Integer);
          $current = 1;
          continue syncloop;

        case 3:
          $t.fastbox(5678, $g.____testlib.basictypes.Integer);
          return;

        default:
          return;
      }
    }
  };
});
