$module('castrejection', function () {
  var $static = this;
  $static.TEST = function () {
    var a;
    var b;
    var somevalue;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          somevalue = $t.fastbox('hello', $g.____testlib.basictypes.String);
          try {
            var $expr = $t.cast(somevalue, $g.____testlib.basictypes.Integer, false);
            a = $expr;
            b = null;
          } catch ($rejected) {
            b = $rejected;
            a = null;
          }
          $current = 1;
          continue syncloop;

        case 1:
          return $t.fastbox((a == null) && !(b == null), $g.____testlib.basictypes.Boolean);

        default:
          return;
      }
    }
  };
});
