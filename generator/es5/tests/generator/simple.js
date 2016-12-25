$module('simple', function () {
  var $static = this;
  $static.SomeGenerator = function () {
    var $current = 0;
    var $continue = function ($yield, $yieldin, $reject, $done) {
      while (true) {
        switch ($current) {
          case 0:
            $yield($t.fastbox(false, $g.____testlib.basictypes.Boolean));
            $current = 1;
            return;

          case 1:
            $yield($t.fastbox(true, $g.____testlib.basictypes.Boolean));
            $current = 2;
            return;

          default:
            $done();
            return;
        }
      }
    };
    return $generator.new($continue, false);
  };
  $static.TEST = $t.markpromising(function () {
    var $result;
    var $temp0;
    var $temp1;
    var v;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            v = null;
            $current = 1;
            $continue($resolve, $reject);
            return;

          case 1:
            $temp1 = $g.simple.SomeGenerator();
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 2:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            value = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              $continue($resolve, $reject);
              return;
            } else {
              $current = 5;
              $continue($resolve, $reject);
              return;
            }
            break;

          case 4:
            v = value;
            $current = 2;
            $continue($resolve, $reject);
            return;

          case 5:
            $resolve(v);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
