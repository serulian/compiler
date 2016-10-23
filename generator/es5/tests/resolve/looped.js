$module('looped', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var $temp0;
    var $temp1;
    var casted;
    var index;
    var value;
    var values;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.List($t.struct).forArray([$t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox(3, $g.____testlib.basictypes.Integer)]).then(function ($result0) {
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
            $g.____testlib.basictypes.Integer.$range($t.fastbox(0, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
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
            index = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 6;
              continue;
            } else {
              $current = 11;
              continue;
            }
            break;

          case 6:
            values.$index(index).then(function ($result0) {
              $result = $result0;
              $current = 7;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 7:
            value = $result;
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
            if (casted != null) {
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
            $resolve($t.fastbox(false, $g.____testlib.basictypes.Boolean));
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
