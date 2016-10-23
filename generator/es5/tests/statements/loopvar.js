$module('loopvar', function () {
  var $static = this;
  this.$class('8caa43e9', 'SomeStream', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.wasChecked = $t.fastbox(false, $g.____testlib.basictypes.Boolean);
      return $promise.resolve(instance);
    };
    $instance.Next = function () {
      var $this = this;
      var $result;
      var r;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              r = $this.wasChecked;
              $this.wasChecked = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
              $g.____testlib.basictypes.Tuple($g.____testlib.basictypes.Boolean, $g.____testlib.basictypes.Boolean).Build($t.fastbox(true, $g.____testlib.basictypes.Boolean), $t.fastbox(!r.$wrapped, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
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
              $resolve($result);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Next|2|29dc432d<4499960a<5ab5941e,5ab5941e>>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomething = function (somethingElse) {
    var $result;
    var $temp0;
    var $temp1;
    var something;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $t.fastbox(1234, $g.____testlib.basictypes.Integer);
            $current = 1;
            continue;

          case 1:
            $temp1 = somethingElse;
            $current = 2;
            continue;

          case 2:
            $temp1.Next().then(function ($result0) {
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
            something = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              continue;
            } else {
              $current = 5;
              continue;
            }
            break;

          case 4:
            $t.fastbox(7654, $g.____testlib.basictypes.Integer);
            $current = 2;
            continue;

          case 5:
            $t.fastbox(5678, $g.____testlib.basictypes.Integer);
            $resolve();
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $temp0;
    var $temp1;
    var i;
    var result;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            result = $t.fastbox('noloop', $g.____testlib.basictypes.String);
            $g.loopvar.SomeStream.new().then(function ($result0) {
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
            s = $result;
            $current = 2;
            continue;

          case 2:
            $temp1 = s;
            $current = 3;
            continue;

          case 3:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
            i = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 5;
              continue;
            } else {
              $current = 6;
              continue;
            }
            break;

          case 5:
            result = i;
            $current = 3;
            continue;

          case 6:
            $resolve(result);
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
