$module('looped', function () {
  var $static = this;
  $static.TEST = function () {
    var $temp0;
    var $temp1;
    var casted;
    var index;
    var values;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.List($t.any).forArray([$t.box(1, $g.____testlib.basictypes.Integer), $t.box(true, $g.____testlib.basictypes.Boolean), $t.box(3, $g.____testlib.basictypes.Integer)]).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            values = $result;
            $current = 2;
            continue;

          case 2:
            $g.____testlib.basictypes.Integer.$range($t.box(0, $g.____testlib.basictypes.Integer), $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $temp1 = $result;
            $current = 4;
            continue;

          case 4:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 5;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 5:
            $result;
            index = $temp0.First;
            if ($t.unbox($temp0.Second)) {
              $current = 6;
              continue;
            } else {
              $current = 11;
              continue;
            }
            break;

          case 6:
            values.$index(index).then(function ($result0) {
              value = $result0;
              $result = value;
              $current = 7;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 7:
            $result;
            $promise.new(function ($resolve) {
              $resolve($t.cast(value, $g.____testlib.basictypes.Boolean, false));
            }).then(function ($result0) {
              casted = $result0;
              $current = 8;
              $continue($resolve, $reject);
              return;
            }).catch(function ($rejected) {
              casted = null;
              $current = 8;
              $continue($resolve, $reject);
              return;
            });
            return;

          case 8:
            if (!(casted == null)) {
              $current = 9;
              continue;
            } else {
              $current = 10;
              continue;
            }
            break;

          case 9:
            $resolve(casted);
            return;

          case 10:
            $current = 4;
            continue;

          case 11:
            $resolve($t.box(false, $g.____testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});