$module('asyncstruct', function () {
  var $static = this;
  this.$struct('4e1a7940', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Foo, Bar) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        Foo: Foo,
        Bar: Bar,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Foo', 'Foo', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'Bar', 'Bar', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<4e1a7940>": true,
        "equals|4|29dc432d<43834c3f>": true,
        "Stringify|2|29dc432d<5cffd9b5>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<4e1a7940>": true,
        "String|2|29dc432d<5cffd9b5>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('dd7aa26e', function (s) {
    return $g.____testlib.basictypes.String.$plus(s.Foo.String(), s.Bar.String());
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var vle;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.translate($g.asyncstruct.DoSomethingAsync($g.asyncstruct.SomeStruct.new($t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer)))).then(function ($result0) {
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
            vle = $result;
            $resolve($g.____testlib.basictypes.String.$equals(vle, $t.fastbox("12", $g.____testlib.basictypes.String)));
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
