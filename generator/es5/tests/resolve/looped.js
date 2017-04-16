$module('looped', function () {
  var $static = this;
  $static.TEST = function () {
    var $temp0;
    var $temp1;
    var casted;
    var index;
    var value;
    var values;
    var $current = 0;
    syncloop: while (true) {
      switch ($current) {
        case 0:
          values = $g.____testlib.basictypes.Slice($t.struct).overArray([$t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox(3, $g.____testlib.basictypes.Integer)]);
          $current = 1;
          continue syncloop;

        case 1:
          $temp1 = $g.____testlib.basictypes.Integer.$range($t.fastbox(0, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer));
          $current = 2;
          continue syncloop;

        case 2:
          $temp0 = $temp1.Next();
          index = $temp0.First;
          if ($temp0.Second.$wrapped) {
            $current = 3;
            continue syncloop;
          } else {
            $current = 7;
            continue syncloop;
          }
          break;

        case 3:
          value = values.$index(index);
          try {
            var $expr = $t.cast(value, $g.____testlib.basictypes.Boolean, false);
            casted = $expr;
          } catch ($rejected) {
            casted = null;
          }
          $current = 4;
          continue syncloop;

        case 4:
          if (casted != null) {
            $current = 5;
            continue syncloop;
          } else {
            $current = 6;
            continue syncloop;
          }
          break;

        case 5:
          return casted;

        case 6:
          $current = 2;
          continue syncloop;

        case 7:
          return $t.fastbox(false, $g.____testlib.basictypes.Boolean);

        default:
          return;
      }
    }
  };
});
