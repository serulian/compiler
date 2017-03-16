$module('resolvereject', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var b;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          try {
            var $expr = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
            a = $expr;
            b = null;
          } catch ($rejected) {
            b = $t.ensureerror($rejected);
            a = null;
          }
          $current = 1;
          continue syncloop;

        case 1:
          return a;

        default:
          return;
      }
    }
  };
});
